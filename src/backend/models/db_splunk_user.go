package models

import (
	"bufio"
	"errors"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SaveSplunkUser(splunkUser SplunkUser) error {
	var conf = settings.GetConfig()
	var searchHeadNames []string
	if splunkUser.Env == "prod" && splunkUser.SearchHead == conf.ProdSHC {
		searchHeads, err := GetSplunkHosts([]string{splunkUser.Env, "SHCMember"})
		if err != nil {
			return err
		}
		for _, sh := range searchHeads {
			searchHeadNames = append(searchHeadNames, sh.Name)
		}
	} else {
		searchHeadNames = append(searchHeadNames, splunkUser.SearchHead)
	}
	if len(searchHeadNames) < 1 {
		return errors.New("No Splunk host found for user: " + splunkUser.UserName + ", env: " + splunkUser.Env)
	}

	err := SpkCreateSplunkApp(splunkUser.UserName, searchHeadNames)
	if err != nil {
		return err
	}

	err = SpkCreateSplunkRole(splunkUser.UserName, searchHeadNames)
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)

	if splunkUser.Env == "prod" {
		splunkUser.PositionIDs = strings.ToLower(splunkUser.PositionIDs)
		SpkCreateOrUpdateSplunkSamlGroup(splunkUser.UserName, splunkUser.SearchHead, splunkUser.PositionIDs)
	}

	SpkCreateSplunkUser(splunkUser.UserName, splunkUser.Password, splunkUser.Email, splunkUser.PositionIDs, searchHeadNames)
	res := mysqldb.Save(&splunkUser)
	if res.Error != nil {
		return res.Error
	}

	go func() {
		SendEmail(conf.EmailSender, conf.AdminEmail, "common", "admin", "create-splunk-account",
			"[SPaaS] Created Splunk Web Account", conf.AppUrl, splunkUser, "")
		SendEmail(conf.EmailSender, splunkUser.Email, "common", splunkUser.UserName, "create-splunk-account",
			"[SPaaS] Created Splunk Web Account", conf.AppUrl, splunkUser, "")
	}()

	return nil
}

func GetSplunkUsers() ([]SplunkUser, error) {
	var splunkUsers []SplunkUser
	res := mysqldb.Order("id").Find(&splunkUsers)
	if res.Error != nil {
		return nil, res.Error
	}
	return splunkUsers, nil
}

func GetSplunkUser(splunkUserName string, env string) (SplunkUser, error) {
	var splunkUser SplunkUser
	res := mysqldb.Where("user_name = ? AND env like ?", splunkUserName, env).Find(&splunkUser)
	if res.Error != nil { //record not found
		return SplunkUser{}, res.Error
	}
	return splunkUser, nil
}

func UpdateSplunkUser(splunkUserName string, newSplunkUser SplunkUser) error {
	var splunkUser SplunkUser
	var oldSplunkUser SplunkUser
	res := mysqldb.Where("user_name = ? ", splunkUserName).Find(&splunkUser)
	if res.Error != nil { //record not found
		return res.Error
	}
	oldSplunkUser = splunkUser
	splunkUser.PositionIDs = strings.ToLower(newSplunkUser.PositionIDs)
	splunkUser.Memo = newSplunkUser.Memo
	splunkUser.RpaasUserName = newSplunkUser.RpaasUserName
	res = mysqldb.Save(&splunkUser)
	if res.Error != nil {
		return res.Error
	}
	if splunkUser.Env == "prod" {
		SpkCreateOrUpdateSplunkSamlGroup(splunkUser.UserName, splunkUser.SearchHead, splunkUser.PositionIDs)
	}

	go func() {
		var conf = settings.GetConfig()
		SendEmail(conf.EmailSender, conf.AdminEmail, "common", "admin", "update-splunk-account",
			"[SPaaS] Updated Splunk Web Account", conf.AppUrl, splunkUser, oldSplunkUser)
		SendEmail(conf.EmailSender, splunkUser.Email, "common", splunkUser.UserName, "update-splunk-account",
			"[SPaaS] Updated Splunk Web Account", conf.AppUrl, splunkUser, oldSplunkUser)
	}()

	return nil
}

func DeleteSplunkUser(splunkUserName string) error {
	splunkUser, err := GetSplunkUser(splunkUserName, "")
	if err != nil {
		return err
	}
	return mysqldb.Delete(&splunkUser).Error
}

func SpkCreateSplunkApp(userName string, hosts []string) error {
	values := url.Values{}
	appName := userName[5:]
	values.Add("name", appName)
	for _, host := range hosts {
		resp, _ := GetSplunkHttpResponse(values, "POST", host, "apps", "")
		defer resp.Body.Close()
		settings.WriteInfoLog("SpkCreateSplunkApp: " + resp.Status)
		if resp.StatusCode != 201 && resp.StatusCode != 409 { // or 409
			return errors.New("Failed to create app: " + appName)
		} else if resp.StatusCode == 409 {
			settings.WriteInfoLog("App exists: " + appName)
		}
	}
	return nil
}

func SpkCreateSplunkRole(userName string, hosts []string) error {
	values := url.Values{}
	values.Add("name", userName)
	values.Add("defaultApp", userName[5:])
	values.Add("srchFilter", "tag=tag_"+userName[5:])
	values.Add("srchIndexesDefault", "*")
	values.Add("srchIndexesAllowed", "*")
	values.Add("srchIndexesAllowed", "_*")
	for _, host := range hosts {
		resp, _ := GetSplunkHttpResponse(values, "POST", host, "roles", "")
		defer resp.Body.Close()
		settings.WriteInfoLog("SpkCreateSplunkRole: " + resp.Status)
		if resp.StatusCode != 201 && resp.StatusCode != 409 { // or 409
			return errors.New("Failed to create role: " + userName)
		} else if resp.StatusCode == 409 {
			settings.WriteInfoLog("Role exists: " + userName)
		}
	}
	return nil
}

func SpkCreateSplunkUser(userName, password, email, positionID string, hosts []string) error { // error will be ignored
	values := url.Values{}
	values.Add("name", userName)
	values.Add("password", password)
	values.Add("defaultApp", userName[5:])
	values.Add("roles", userName)
	values.Add("roles", "user_without_rt_search")
	values.Add("email", email)
	for _, host := range hosts {
		resp, _ := GetSplunkHttpResponse(values, "POST", host, "users", "")
		defer resp.Body.Close()
		settings.WriteInfoLog("SpkCreateSplunkUser: " + resp.Status)
		if resp.StatusCode != 201 {
			//send email
			return errors.New("Failed to create user: " + userName) // returns 400 if user exists
		}
	}
	return nil
}

func SpkCreateOrUpdateSplunkSamlGroup(userName, host, positionIDs string) error { // error will be ignored
	IDs := strings.Split(positionIDs, ",")
	settings.WriteInfoLog(IDs)
	if len(IDs) < 1 {
		return errors.New("No position ID specified: " + userName)
	}
	for _, ID := range IDs {
		if ID == "" {
			continue
		}
		group := SpkGetSamlRoles("rogin-"+ID, host)
		roleStr := ""
		values := url.Values{}
		values.Add("roles", "user_without_rt_search")
		values.Add("roles", userName)
		for _, role := range group.Roles {
			roleStr = roleStr + role + ","
			values.Add("roles", role)
		}
		settings.WriteInfoLog("Current roles in group " + ID + ": " + roleStr)

		resp, _ := GetSplunkHttpResponse(values, "POST", host, "saml_groups", "/rogin-"+ID)
		defer resp.Body.Close()
		if resp.StatusCode != 200 { // or 400
			// send email
			settings.WriteInfoLog("Failed to create/update saml group: " + "rogin-" + ID)
			// return errors.New("Failed to create/update saml group: " + "rogin-" + ID)
		}
	}
	return nil
}

func SpkGetSamlRoles(groupName, host string) SamlGroup {
	var group SamlGroup
	resp, _ := GetSplunkHttpResponse(url.Values{}, "GET", host, "saml_groups", "/"+groupName)
	defer resp.Body.Close()
	settings.WriteInfoLog("SpkGetSamlRoles: " + resp.Status)
	if resp.StatusCode != 200 { // 404 Not Found
		settings.WriteInfoLog("Saml group not found, creating saml group: " + groupName)
		return SamlGroup{}
	}
	group, err := parseSamlGroupResponse(resp)
	if err != nil {
		settings.WriteInfoLog("Saml group not found, creating saml group: " + groupName)
		return SamlGroup{}
	}
	group.Name = groupName
	return group
}

func parseSamlGroupResponse(resp *http.Response) (SamlGroup, error) {
	var group SamlGroup
	ch := make(chan string)
	groupFound := false
	adminCount := 0
	go func(ch chan string) {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			if strings.Contains(string(line), "<s:item>") { // e.g : <s:item>user_without_rt_search</s:item>
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
				end := strings.Index(line, "</s:item>")
				if start < 0 || end < 0 || start+1 > end {
					continue
				}
				role := line[start+1 : end]
				if strings.HasPrefix(role, "user_") || strings.HasPrefix(role, "func_") {
					group.Roles = append(group.Roles, role)
					groupFound = true
				}
				if strings.HasPrefix(role, "admin") {
					adminCount++
					if adminCount == 3 {
						group.Roles = append(group.Roles, role)
						groupFound = true
					}
				}
			}
		case <-time.After(1 * time.Second):
		}
	}
	if !groupFound {
		return SamlGroup{}, errors.New("Saml group not found")
	}
	return group, nil
}
