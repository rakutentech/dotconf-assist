package models

import (
	"bufio"
	"bytes"
	"errors"
	// "fmt"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"golang.org/x/crypto/ssh"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
1. Update deployment (bind app with serverclass)
2. Apply deployment
	2.1 CreateDeploymentApp (create dir on deployment server), status: requested
	2.2 CreateDeploymentServerClass (, )
		2.2.1 create server class on deployment server (serverclasses.conf)
		2.2.2 Update distributed app (serverclasses.conf), deployment client will restart after this step, status: deployed
	2.3 CreateTag
		2.3.1 fetch fields for current tag(name_<name>)
		2.3.2 calc fields need to be added to current tag
		2.3.3 change permission for fields, host=host*
		2.3.4 verify fileds permission
3. Deinstall deployment app
	3.1 Deinstall app in Splunk
	3.2 remove dir created in 2.1
*/

type DeploymentInfo struct {
	User     string
	Env      string
	App      string
	Success  string
	Error    string
	Link4Tag string
}

func SaveDeployment(deploymentList []Deployment) error { // not used
	for _, dply := range deploymentList {
		res := mysqldb.Save(&dply)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func CreateDeploymentApp(deploymentApp DeploymentApp) error {
	// UpdateAppFieldsByID(deploymentApp.AppID, []string{"deploy_status"}, []interface{}{2}, "apps") // requested
	var conf settings.Configuration
	conf = settings.GetConfig()
	config := &ssh.ClientConfig{
		User: conf.DPLYSSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.DPLYSSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	baseDir := "/opt/splunk/etc/deployment-apps"
	unixAppTemplateDir := "/opt/splunk/etc/deployment-apps/Splunk_TA_nix_basic"
	// baseDir := "splunk-test/etc/deployment-apps"
	// unixAppTemplateDir := "splunk-test/etc/deployment-apps/Splunk_TA_nix_basic"
	appDir := baseDir + "/" + deploymentApp.FolderName
	binDir := appDir + "/" + "bin"
	localDir := appDir + "/" + "local"
	inputsConfPath := localDir + "/" + "inputs.conf"
	var cmds []string
	if deploymentApp.AppType == "app" {
		cmds = append(cmds, "rm -rf "+appDir)
		cmds = append(cmds, "mkdir -p "+localDir)
		cmds = append(cmds, "chmod 777 -R "+appDir)
		cmds = append(cmds, "echo \""+deploymentApp.InputsConf+"\" > "+inputsConfPath)
		cmds = append(cmds, "cat "+inputsConfPath)
		if len(deploymentApp.ScriptIDs) > 0 {
			cmds = append(cmds, "mkdir -p "+binDir)
			for _, inputID := range deploymentApp.ScriptIDs {
				var anyInput interface{}
				var err error
				if anyInput, err = GetInput(inputID.ID, "script"); err != nil {
					return err
				}
				var scriptInput ScriptInput
				scriptInput = anyInput.(ScriptInput)
				scriptInput.ScriptCode = (string)(scriptInput.Script)
				scriptInput.ScriptCode = strings.Replace(scriptInput.ScriptCode, "\"", "\\\"", -1)
				// settings.WriteInfoLog(scriptInput.ScriptCode)
				cmds = append(cmds, "echo \""+scriptInput.ScriptCode+"\" >"+binDir+"/"+scriptInput.ScriptName)
			}
			cmds = append(cmds, "chmod 777 -R "+binDir)
		}
	} else { //unix app // copy folder and update intputs.conf
		cmds = append(cmds, "rm -rf "+appDir)
		cmds = append(cmds, "mkdir -p "+appDir)
		cmds = append(cmds, "chmod 777 -R "+appDir)
		cmds = append(cmds, "cp -pr "+unixAppTemplateDir+"/* "+appDir)
		cmds = append(cmds, "echo \""+deploymentApp.InputsConf+"\" > "+inputsConfPath)
		cmds = append(cmds, "cat "+inputsConfPath)
	}
	var splunkHost SplunkHost
	res := mysqldb.Where("role = ? AND env = ?", "DPLY", deploymentApp.Env).Find(&splunkHost)
	if res.Error != nil { //record not found
		return res.Error
	}
	err := executeCmd(cmds, splunkHost.Name, config)
	if err != nil {
		return err
	} else {
		settings.WriteInfoLog("Created deployment app: " + deploymentApp.FolderName)
	}
	return nil
}

func CreateDeploymentServerClass(dplySCs []DeploymentServerClass) error { // either create or update
	var result4UserOK string
	var result4UserError string
	var result4AdminOK string
	var result4AdminError string
	var link4Tag string
	var user User
	var err error
	var success bool
	success = true
	user, err = GetUser(dplySCs[0].User)
	var splunkUser SplunkUser
	splunkUser, err = GetSplunkUser(dplySCs[0].User, dplySCs[0].Env)

	if len(dplySCs) < 1 {
		return errors.New("No server class specified")
	}

	var splunkHost SplunkHost
	res := mysqldb.Where("role = ? AND env = ?", "DPLY", dplySCs[0].Env).Find(&splunkHost)
	if res.Error != nil { //record not found
		return res.Error
	}

	for i, sc := range dplySCs {
		err = SpkCheckServerClassExistence(sc.Env, sc.ServerClassName)
		if err != nil { // sc not exists, need to create
			settings.WriteInfoLog("Server class doesn't exist, creating server class: " + sc.ServerClassName)
			err = SpkCreateDeploymentServerClass(sc, splunkHost.Name)
			if err != nil {
				success = false
				// return errors.New("Failed to create server class: " + err.Error())
			}
		} else { // update sc
			settings.WriteInfoLog("Server class exists, updating server class: " + sc.ServerClassName)
			err = SpkUpdateDeploymentServerClass(sc, splunkHost.Name)
			if err != nil {
				success = false
				// return errors.New("Failed to update server class: " + err.Error())
			}
		}

		err = SpkUpdateDistributedApp(sc.ServerClassName, sc.AppName, splunkHost.Name)
		if err != nil {
			success = false
			// return errors.New("Failed to update distributed app: " + err.Error())
		}
		dplySCs[i].ServerClassName = dplySCs[i].ServerClassName[len("sc_"+sc.User+"_"):] // for email
	}

	var subject string
	if success {
		result4UserOK = "It has been deployed sucessfully, there is no further action required."
		result4AdminOK = "It has been deployed sucessfully."
		err = CreateTag(dplySCs[0].User, dplySCs[0].Env, dplySCs)
		if err != nil {
			settings.WriteInfoLog("Failed to create tag: " + err.Error())
			link4Tag = "https://" + splunkUser.SearchHead + "/en-US/manager/search/admin/tags"
			result4AdminError = err.Error() + " (tag_" + dplySCs[0].User[5:] + ")"
			subject = "[SPaaS] [Need Action] Request Deploy Setting"
		} else {
			result4AdminOK += " There is no further action required."
			subject = "[SPaaS] Request Deploy Setting"
		}
		settings.WriteInfoLog("Created deployment server classes")
	} else {
		result4AdminError = "It failed, please check the details for user."
		result4UserError = "It failed, please check your server class name and try to apply again later."
		subject = "[SPaaS] [Need Action] Request Deploy Setting"
	}

	go func() {
		var conf settings.Configuration
		conf = settings.GetConfig()
		err = SendEmail(conf.EmailSender, conf.AdminEmail, "common", "admin", "deploy-app",
			subject, conf.AppUrl, dplySCs,
			DeploymentInfo{dplySCs[0].User, dplySCs[0].Env,
				dplySCs[0].AppName[len("dsapp_"+dplySCs[0].User+"_"):], result4AdminOK, result4AdminError, link4Tag})
		if err != nil {
			settings.WriteInfoLog(err)
		}

		err = SendEmail(conf.EmailSender, user.Email, "common", dplySCs[0].User, "deploy-app",
			"[SPaaS] Request Deploy Setting", conf.AppUrl, dplySCs,
			DeploymentInfo{dplySCs[0].User, dplySCs[0].Env,
				dplySCs[0].AppName[len("dsapp_"+dplySCs[0].User+"_"):], result4UserOK, result4UserError, ""})
		if err != nil {
			settings.WriteInfoLog(err)
		}
	}()
	if success { // send email to user only succeeded
		return UpdateAppFieldsByID(dplySCs[0].AppID, []string{"deploy_status"}, []interface{}{1}, "apps") // configured
	} else {
		return errors.New("Error occured when creating server class")
	}
}

func CreateTag(user, env string, dplySCs []DeploymentServerClass) error {
	var splunkUser SplunkUser
	var tag Tag
	res := mysqldb.Where("user_name = ? AND env = ?", user, env).Find(&splunkUser)
	if res.Error != nil { //record not found
		return res.Error
	}

	tagName := "tag_" + user[5:]
	tag, err := SpkGetTagFields(tagName, splunkUser.SearchHead)
	if err != nil {
		settings.WriteInfoLog(err.Error())
	}

	var fwdrs []string
	var fwdrStr string
	var existingFwdrStr string
	var fwdr2AddStr string
	for _, sc := range dplySCs {
		for _, fwdr := range sc.ForwarderNames {
			fwdrs = append(fwdrs, fwdr.Name)
			fwdrStr = fwdrStr + fwdr.Name + ","
		}
	}
	for _, fd := range tag.Hosts {
		existingFwdrStr = existingFwdrStr + fd + ","
	}

	settings.WriteInfoLog("Forwarders in server class: " + fwdrStr)
	settings.WriteInfoLog("Existing forwarders for tag " + tagName + ": " + existingFwdrStr)
	fwdrWildcards := calcForwarderWildcards(fwdrs, tag.Hosts)
	settings.WriteInfoLog("Forwarders to be added: " + fwdr2AddStr)

	if len(fwdrWildcards) < 1 {
		return nil
	}

	err = SpkCreateTag(tagName, splunkUser.SearchHead, fwdrWildcards)
	if err != nil {
		return err
	}

	fwdrWithProblem := SpkChangeTagPermission(tagName, splunkUser.SearchHead, fwdrWildcards)
	var fwdrsWithError string
	for _, fwdr := range fwdrWithProblem {
		fwdrsWithError = fwdrsWithError + fwdr + ","
	}
	if len(fwdrWithProblem) > 0 {
		settings.WriteInfoLog("Need to manually change permision for: " + fwdrsWithError)
		return errors.New("Need to manually change permision for: " + fwdrsWithError)
	}
	//<msg type="ERROR">Cannot overwrite existing app object</msg>
	//link https://hostname/en-US/manager/search/admin/tags

	fwdrsWithError = ""
	fwdrWithProblem = SpkVerifyTagPermission(tagName, splunkUser.SearchHead, fwdrWildcards)
	for _, fwdr := range fwdrWithProblem {
		fwdrsWithError = fwdrsWithError + fwdr + ","
	}
	if len(fwdrWithProblem) > 0 {
		settings.WriteInfoLog("Need to manually change permision for: " + fwdrsWithError)
		return errors.New("Need to manually change permision for: " + fwdrsWithError)
	}
	return nil
}

func RemoveDeploymentServerClass(dplySCs []DeploymentServerClass) error { // not used. not that important, returned value will not be used
	if len(dplySCs) < 1 {
		settings.WriteInfoLog("Failed to delete server class: No server class specified")
		return errors.New("No server class specified")
	}
	var splunkHost SplunkHost
	res := mysqldb.Where("role = ? AND env = ?", "DPLY", dplySCs[0].Env).Find(&splunkHost)
	if res.Error != nil { //record not found
		settings.WriteInfoLog("Failed to delete server class: " + res.Error.Error())
		return errors.New("Failed to update server class: " + res.Error.Error())
	}

	for _, sc := range dplySCs {
		err := SpkDeleteDeploymentServerClass(sc, splunkHost.Name)
		if err != nil {
			settings.WriteInfoLog(err.Error())
			return errors.New(err.Error())
		}
	}

	settings.WriteInfoLog("Remove deployment server class successfully")
	return nil
}

func DeinstallDeploymentApp(deploymentApp DeploymentApp) error {
	var splunkHost SplunkHost
	res := mysqldb.Where("role = ? AND env = ?", "DPLY", deploymentApp.Env).Find(&splunkHost)
	if res.Error != nil { //record not found
		return res.Error
	}

	err := SpkDeinstallDeploymentApp(deploymentApp.FolderName, splunkHost.Name)
	if err != nil {
		return errors.New("Failed to deinstall distributed app: " + err.Error())
	}

	// err = RemoveDeploymentServerClass()
	// if err != nil {
	// 	settings.WriteInfoLog("Deinstalled the app, but failed to remove server class from deployment server: " + err.Error())
	// }

	settings.WriteInfoLog("Deinstalled app " + deploymentApp.FolderName + " successfully")
	err = RemoveDeploymentApp(deploymentApp)
	if err != nil {
		settings.WriteInfoLog("Deinstalled the app, removed server class, but failed to delete the folder from deployment server: " + err.Error())
		return nil
		// return errors.New("Failed to deinstall distributed app: " + err.Error())
	}

	return UpdateAppFieldsByID(deploymentApp.AppID, []string{"deploy_status"}, []interface{}{3}, "apps")
}

func RemoveDeploymentApp(deploymentApp DeploymentApp) error { // remove app from deployment server under deployment-apps
	var conf settings.Configuration
	conf = settings.GetConfig()
	config := &ssh.ClientConfig{
		User: conf.DPLYSSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.DPLYSSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	baseDir := "/opt/splunk/etc/deployment-apps"
	// baseDir := "splunk-test/etc/deployment-apps"
	appDir := baseDir + "/" + deploymentApp.FolderName
	// unixAppTemplateDir := "splunk-test/etc/deployment-apps/Splunk_TA_nix_basic"
	var cmds []string
	cmds = append(cmds, "rm -rf "+appDir)
	var splunkHost SplunkHost
	res := mysqldb.Where("role = ? AND env = ?", "DPLY", deploymentApp.Env).Find(&splunkHost)
	if res.Error != nil { //record not found
		return res.Error
	}
	err := executeCmd(cmds, splunkHost.Name, config)
	if err != nil {
		return err
	} else {
		settings.WriteInfoLog("Deleted folder from deployment server: " + appDir)
	}
	return nil
}

func GetDeploymentList(envUserFrom []string) ([]Deployment, error) {
	var deploymentList []Deployment
	res := mysqldb.Where("env LIKE ? AND user_name LIKE ?", envUserFrom[0], envUserFrom[1]).Order("id").Find(&deploymentList)
	if res.Error != nil { // res.Error is nil even if no record found
		return nil, res.Error
	}

	return deploymentList, nil
}

func GetServerClassIDsByAppID(appID int, getForwarders bool) ([]int, error) {
	var deploymentList []Deployment
	var scIDs []int
	res := mysqldb.Select("server_class_id").Where("app_id = ?", appID).Find(&deploymentList)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, dply := range deploymentList {
		scIDs = append(scIDs, dply.ServerClassID)
	}
	return scIDs, nil
}

func GetDeployment(name, user string) (Deployment, error) {
	var deployment Deployment
	res := mysqldb.Where("name = ? AND user_name = ?", name, user).Find(&deployment)
	if res.Error != nil { //record not found
		return Deployment{}, res.Error
	}
	return deployment, nil
}

func UpdateDeployment(user string, newDplyList []Deployment) error {
	if len(newDplyList) < 1 {
		return nil
	}
	if err := DeleteDeploymentByAppID(newDplyList[0].AppID); err != nil {
		return err
	}
	for _, dply := range newDplyList {
		res := mysqldb.Save(&dply)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func DeleteDeploymentByIDs(appID int, scIDs []int) error {
	res := mysqldb.Where("app_id = ? AND server_class_id IN (?)", appID, scIDs).Delete(Deployment{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteDeploymentByAppID(appID int) error {
	res := mysqldb.Where("app_id = ?", appID).Delete(Deployment{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteDeployment(name, user string) error {
	deployment, err := GetDeployment(name, user)
	if err != nil {
		return err
	}
	return mysqldb.Delete(&deployment).Error
}

func executeCmd(cmd []string, hostname string, config *ssh.ClientConfig) error {
	conn, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		return errors.New("Failed to connect:" + err.Error())
	}
	defer conn.Close()
	for _, command := range cmd {
		var stdoutBuf bytes.Buffer
		session, err := conn.NewSession()
		if err != nil {
			return errors.New("Failed to create session: " + err.Error())
		}
		defer session.Close()
		session.Stdout = &stdoutBuf
		err = session.Run(command)
		if err != nil {
			return errors.New("Failed to run command: '" + command + "': " + err.Error())
		}
		// fmt.Println(stdoutBuf.String())
	}
	// return stdoutBuf.String()
	return nil
}

func SpkCheckServerClassExistence(env, sc string) error {
	var splunkHost SplunkHost
	res := mysqldb.Where("role = ? AND env = ?", "DPLY", env).Find(&splunkHost)
	if res.Error != nil { //record not found
		return res.Error
	}

	resp, _ := GetSplunkHttpResponse(url.Values{}, "GET", splunkHost.Name, "server_classes", "/"+sc)
	settings.WriteInfoLog("SpkCheckServerClassExistence: " + resp.Status)
	if resp.StatusCode != 200 {
		return errors.New("Server class not found: " + sc)
	}
	defer resp.Body.Close()
	return nil
}

func SpkGetTagFields(tagName, host string) (Tag, error) {
	var tag Tag
	resp, _ := GetSplunkHttpResponse(url.Values{}, "GET", host, "tags", "/"+tagName)
	defer resp.Body.Close()
	settings.WriteInfoLog("SpkGetTagFields: " + resp.Status)
	if resp.StatusCode != 200 { // 404 Not Found
		return Tag{}, errors.New("Tag not found: " + tagName)
	}
	tag, err := parseTagResponse(resp)
	if err != nil {
		return Tag{}, err
	}
	return tag, nil
}

func SpkCreateDeploymentServerClass(sc DeploymentServerClass, host string) error {
	values := url.Values{}
	values.Add("name", sc.ServerClassName)
	values.Add("restartSplunkd", "true")
	for i, fwdr := range sc.ForwarderNames {
		values.Add("whitelist."+strconv.Itoa(i), fwdr.Name)
	}
	resp, _ := GetSplunkHttpResponse(values, "POST", host, "server_classes", "")
	defer resp.Body.Close()
	settings.WriteInfoLog("SpkCreateDeploymentServerClass: " + sc.ServerClassName + "Status: " + resp.Status)
	if resp.StatusCode != 201 {
		return errors.New("Failed to create server class: " + sc.ServerClassName)
	}
	return nil
}

//e.g forwarders: dm-planet101, dm-planet102
// existingWildcards: dm-planet*,
func calcForwarderWildcards(forwarders []string, existingWildcards []string) []string {
	var wildcards []string
	var wildcards2Add []string
	re := regexp.MustCompile("[0-9]")
	for _, fwdr := range forwarders {
		loc := re.FindStringIndex(fwdr)
		if len(loc) > 0 {
			wildcards = append(wildcards, fwdr[0:loc[0]]+"*")
		} else { //no digit found
			wildcards = append(wildcards, fwdr)
		}
	}

	wildcards = unique(wildcards)

	for _, wc := range wildcards {
		ignore := false
		for _, existwc := range existingWildcards {
			if existwc[len(existwc):] == "*" { //last char is *
				if strings.HasPrefix(wc, existwc[0:len(existwc)-1]) { //nsmon*, ns
					ignore = true
					break
				}
			} else { // last char is not *
				if strings.HasPrefix(wc, existwc) { //nsmon*, test.local
					ignore = true
					break
				}
			}
		}
		if !ignore {
			wildcards2Add = append(wildcards2Add, wc)
		}
	}

	return wildcards2Add
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func SpkCreateTag(tagName, host string, forwarderWildcards []string) error {
	if len(forwarderWildcards) < 1 {
		return errors.New("Failed to create tag, no forwarder specified")
	}
	values := url.Values{}
	for _, fwdr := range forwarderWildcards {
		values.Add("add", "host::"+fwdr)
	}

	resp, _ := GetSplunkHttpResponse(values, "POST", host, "tags", "/"+tagName)
	defer resp.Body.Close()
	settings.WriteInfoLog("SpkCreateTag: " + resp.Status) //SpkCreateTag: 201 Created or 200 OK
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return errors.New("Failed to create tag: " + tagName)
	}
	return nil
}

func SpkChangeTagPermission(tagName, host string, forwarderWildcards []string) []string {
	values := url.Values{}
	values.Add("owner", "admin")
	values.Add("sharing", "global")
	values.Add("perms.write", "admin")
	values.Add("perms.read", "*")
	var wildcardsWithProblem []string
	for _, fwdr := range forwarderWildcards {
		resp, _ := GetSplunkHttpResponse(values, "POST", host, "tag_permission", fwdr)
		defer resp.Body.Close()
		settings.WriteInfoLog("SpkChangeTagPermission: " + resp.Status)
		if resp.StatusCode != 200 { // SpkChangeTagPermission: 409 Conflict
			wildcardsWithProblem = append(wildcardsWithProblem, fwdr)
		}
	}
	return wildcardsWithProblem
}

func SpkVerifyTagPermission(tagName, host string, forwarderWildcards []string) []string {
	okKeyword := "name=\"sharing\">global"
	values := url.Values{}
	var wildcardsWithProblem []string
	for _, fwdr := range forwarderWildcards {
		resp, err := GetSplunkHttpResponse(values, "GET", host, "tag_permission", fwdr)
		defer resp.Body.Close()
		settings.WriteInfoLog("SpkVerifyTagPermission: " + resp.Status)
		if err = ParseHttpResponse(resp, okKeyword, "Verify tag permission: "+tagName); err != nil {
			wildcardsWithProblem = append(wildcardsWithProblem, fwdr)
		}
	}
	return wildcardsWithProblem
}

func SpkDeleteDeploymentServerClass(sc DeploymentServerClass, host string) error { // not used
	values := url.Values{}
	resp, _ := GetSplunkHttpResponse(values, "DELETE", host, "server_classes", "/"+sc.ServerClassName)
	defer resp.Body.Close()
	settings.WriteInfoLog("SpkDeleteDeploymentServerClass: " + resp.Status)
	if resp.StatusCode != 200 { //or 500
		return errors.New("Failed to delete deployment server class: " + sc.ServerClassName)
	}
	return nil
}

func SpkUpdateDeploymentServerClass(sc DeploymentServerClass, host string) error {
	values := url.Values{}
	values.Add("restartSplunkd", "true")
	for i, fwdr := range sc.ForwarderNames {
		values.Add("whitelist."+strconv.Itoa(i), fwdr.Name)
	}
	resp, _ := GetSplunkHttpResponse(values, "POST", host, "server_classes", "/"+sc.ServerClassName)
	defer resp.Body.Close()
	settings.WriteInfoLog("SpkUpdateDeploymentServerClass: " + resp.Status)
	if resp.StatusCode != 200 { // or 404
		return errors.New("Failed to update deployment server class: " + sc.ServerClassName)
	}
	return nil
}

func SpkUpdateDistributedApp(scName, appName, host string) error {
	values := url.Values{}
	values.Add("serverclass", scName)
	values.Add("restartSplunkd", "true")
	resp, _ := GetSplunkHttpResponse(values, "POST", host, "distributed_apps", "/"+appName)
	defer resp.Body.Close()
	settings.WriteInfoLog("SpkUpdateDistributedApp: " + resp.Status)
	if resp.StatusCode != 200 { // or 500
		return errors.New("Failed to update destributed app: " + appName)
	}
	return nil
}

func SpkDeinstallDeploymentApp(appName, host string) error {
	okKeyword := "<title>" + appName + "</title>"
	values := url.Values{}
	values.Add("deinstall", "true")
	resp, err := GetSplunkHttpResponse(values, "POST", host, "distributed_apps", "/"+appName)
	defer resp.Body.Close()
	settings.WriteInfoLog("SpkDeinstallDeploymentApp: " + resp.Status)
	// if resp.StatusCode != 200 { // code is always 200 even if app doesn't exist
	// 	return errors.New("Failed to deinstall destributed app: " + appName)
	// }
	if err != nil {
		return err
	}
	return ParseHttpResponse(resp, okKeyword, "deinstall app: "+appName)
}

func parseTagResponse(resp *http.Response) (Tag, error) {
	var tag Tag
	ch := make(chan string)
	tagFound := false
	go func(ch chan string) {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			if strings.Contains(string(line), "<title>") { // e.g : <title>host::dev-dclx*</title>
				ch <- line
			}
		}
	}(ch)

readInLoop:
	for {
		select {
		case line, ok := <-ch:
			if !ok {
				break readInLoop
			} else {
				start := strings.Index(line, ">")
				end := strings.Index(line, "</title>")
				if start < 0 || end < 0 || start+1 > end {
					continue
				}
				field := line[start+1 : end]
				// host::dev-dclx*
				// sourcetype::rpaas_v2
				// source::/as/ab
				// rpaas_guid::abcdedf
				start = strings.Index(field, "::")
				if start < 0 {
					continue
				}
				if strings.Contains(field, "host::") {
					tag.Hosts = append(tag.Hosts, field[start+2:])
				} else if strings.Contains(field, "sourcetype::") {
					tag.SourceTypes = append(tag.SourceTypes, field[start+2:])
				} else if strings.Contains(field, "source::") {
					tag.Sources = append(tag.Sources, field[start+2:])
				} else if strings.Contains(field, "rpaas_guid::") {
					tag.GUIDs = append(tag.GUIDs, field[start+2:])
				}
				tagFound = true
			}
		case <-time.After(1 * time.Second):
		}
	}
	if !tagFound {
		return Tag{}, errors.New("Tag not found")
	}
	return tag, nil
}
