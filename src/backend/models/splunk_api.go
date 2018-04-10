package models

import (
	"bufio"
	"crypto/tls"
	"errors"
	// "fmt"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetSplunkEndpoint(name, param string) string {
	switch name {
	case "users":
		return "/services/authentication/users"
	case "roles":
		return "/services/authorization/roles"
	case "apps":
		return "/services/apps/local"
	case "saml_groups":
		return "/services/admin/SAML-groups" + param
	case "deployment_clients":
		return "/services/deployment/server/clients?count=0"
	case "distributed_apps":
		return "/services/deployment/server/applications" + param
	case "server_classes":
		return "/services/deployment/server/serverclasses" + param
	case "tags":
		return "/servicesNS/admin/search/search/tags" + param
	case "tag_permission":
		return "/servicesNS/admin/search/saved/fvtags/host=" + param + "/acl"
	default:
		return ""
	}
}

func GetSplunkHttpResponse(values url.Values, method, host, resource, param string) (*http.Response, error) {
	var conf = settings.GetConfig()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	proxyString := conf.Proxy
	if proxyString != "" {
		proxyUrl, _ := url.Parse(proxyString)
		tr.Proxy = http.ProxyURL(proxyUrl)
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(20) * time.Second,
	}

	req, err := http.NewRequest(method, "https://"+host+GetSplunkEndpoint(resource, param), strings.NewReader(values.Encode()))
	req.SetBasicAuth(conf.SplunkAdminUsername, conf.SplunkAdminPassword)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func ParseHttpResponse(resp *http.Response, lineKeyWord, msg string) error {
	ch := make(chan string)
	resOK := false
	go func(ch chan string) {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			ch <- line
		}
	}(ch)

readInLoop:
	for {
		select {
		case line, ok := <-ch:
			if !ok {
				break readInLoop
			} else {
				// fmt.Printf("%s\n", line)
				if strings.Contains(line, "type=\"ERROR\">") {
					start := strings.Index(line, ">")
					end := strings.Index(line, "</msg>")
					return errors.New(line[start+1 : end])
				}
				if strings.Contains(line, lineKeyWord) {
					resOK = true
					break
				}
			}
		case <-time.After(1 * time.Second):
		}
	}

	if resOK {
		return nil
	} else {
		return errors.New("Error occured when calling Splunk API: " + msg)
	}
}

//create app

// Convert an external group to internal roles.
// curl -k -u admin:changeme  https://localhost:8089/services/admin/SAML-groups -d name=Splunk -d roles=user

// delete a saml group
// curl -k -u admin:password --request DELETE https://localhost:8089/services/admin/SAML-groups/group_to_delete

// list server class
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/serverclasses

// create new serverclass
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/serverclasses -d name=sc_apps_ombra -d whitelist.1=host1 -d whitelist.2=host2

// remove a server class
// curl -k -u admin:pass --request DELETE https://localhost:8089/services/deployment/server/serverclasses/sc_apps_shadow

// get a server class
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/serverclasses/sc_mach_type

// update a server class
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/serverclasses/sc_apps_ombra -d stateOnClient=noop

// rename a server class
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/serverclasses/rename -d oldName=sc_apps_ombra -d newName=sc_apps_shadow

// List distributed apps
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/applications

// get a distributed app
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/applications/wma-app1

// update a distributed app
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/applications/wma-app3 -d serverclass=sc_apps_wma restartSplunkd=true

// deinstall a destributed app
// curl -k -u admin:pass https://localhost:8089/services/deployment/server/applications/wma-app3 -d deinstall=true

// get tags for a
// curl -k -u admin:pass https://localhost:8089/servicesNS/admin/search/search/tags
