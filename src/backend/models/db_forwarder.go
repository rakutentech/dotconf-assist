package models

import (
	"bufio"
	// "github.com/rakutentech/dotconf-assist/src/backend/settings"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SaveForwarder(forwarders []Forwarder) error {
	for _, fwdr := range forwarders {
		fwdr.CreatedAt = time.Now()
		fwdr.Share = "private"
		res := mysqldb.Save(&fwdr)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func GetForwarders(envUserFrom []string, isAdmin bool) ([]Forwarder, error) {
	var forwarders []Forwarder
	if envUserFrom[2] == "deployment_server" {
		return SpkGetForwarders(envUserFrom[0], envUserFrom[1])
	} else {
		if isAdmin {
			res := mysqldb.Where("env = ?", envUserFrom[0]).Order("id").Find(&forwarders)
			if res.Error != nil { // res.Error is nil even if no record found
				return nil, res.Error
			}
		} else {
			res := mysqldb.Where("env LIKE ? AND user_name LIKE ?", envUserFrom[0], envUserFrom[1]).Order("id").Find(&forwarders)
			if res.Error != nil { // res.Error is nil even if no record found
				return nil, res.Error
			}
		}

	}
	return forwarders, nil
}

func GetForwardersByIDs(fwdrIDs []string) ([]Forwarder, error) {
	var forwarders []Forwarder
	res := mysqldb.Table("forwarders").Where("id IN (?)", fwdrIDs).Find(&forwarders)
	if res.Error != nil { // res.Error is nil even if no record found
		return nil, res.Error
	}
	return forwarders, nil
}

func getForbiddenForwarders(env, user string) ([]Forwarder, error) { //get all forbidden forwarders
	var forwarders []Forwarder
	res := mysqldb.Where("env LIKE ? AND share != ? AND share != ?", env, user, "public").Order("id").Find(&forwarders)
	if res.Error != nil { // res.Error is nil even if no record found
		return nil, res.Error
	}
	return forwarders, nil
}

func GetForwarder(name, user string) (Forwarder, error) {
	var forwarder Forwarder
	res := mysqldb.Where("name = ? AND user_name = ?", name, user).Find(&forwarder)
	if res.Error != nil { //record not found
		return Forwarder{}, res.Error
	}
	return forwarder, nil
}

func UpdateForwarder(name, user string, newForwarder Forwarder) error {
	var forwarder Forwarder
	res := mysqldb.Where("name = ? AND user_name = ?", name, user).Find(&forwarder)
	if res.Error != nil { //record not found
		return res.Error
	}
	// forwarder.Name = newForwarder.Name
	// forwarder.UserName = newForwarder.UserName
	forwarder.Share = newForwarder.Share
	res = mysqldb.Save(&forwarder)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteForwarder(name, user string) error {
	forwarder, err := GetForwarder(name, user)
	if err != nil {
		return err
	}

	DeleteServerClassByForwarderID(forwarder.ID)
	return mysqldb.Delete(&forwarder).Error
}

func DeleteServerClassByForwarderID(forwarderID int) error {
	mysqldb.Where("forwarder_id = ?", forwarderID).Delete(ServerClassForwarder{})
	return nil
}

func SpkGetForwarders(env, user string) ([]Forwarder, error) {
	var forwarders []Forwarder
	var splunkHost SplunkHost

	res := mysqldb.Where("role = ? AND env = ?", "DPLY", env).Find(&splunkHost)
	if res.Error != nil { //record not found
		return nil, res.Error
	}

	resp, err := GetSplunkHttpResponse(url.Values{}, "GET", splunkHost.Name, "deployment_clients", "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	forwarders, err = parseForwarderResponse(resp, env, user)
	if err != nil {
		return nil, err
	}
	return forwarders, nil
}

func parseForwarderResponse(resp *http.Response, env, user string) ([]Forwarder, error) {
	var forwarders []Forwarder
	var forbiddenForwarders []Forwarder
	forbiddenForwarders, err := getForbiddenForwarders(env, user)
	if err != nil { //record not found
		return nil, err
	}

	ch := make(chan string)
	count := 0
	go func(ch chan string) {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			if strings.Contains(string(line), "name=\"hostname\"") { //or dns
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
				end := strings.Index(line, "</s:key>")
				if start < 0 || end < 0 || start+1 > end {
					continue
				}
				host := line[start+1 : end]
				validFwdr := true
				for _, fwdr := range forbiddenForwarders {
					if host == fwdr.Name {
						validFwdr = false
						break
					}
				}

				if validFwdr {
					var forwarder Forwarder
					forwarder.ID = count
					forwarder.Name = host
					forwarders = append(forwarders, forwarder)
					count++
				}
			}
		case <-time.After(1 * time.Second):
		}
	}
	return forwarders, nil
}
