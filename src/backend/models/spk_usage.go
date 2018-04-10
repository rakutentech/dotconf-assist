package models

import (
	"encoding/json"
	"fmt"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"github.com/sebkl/splunk-golang"
	"strconv"
	"time"
)

var splunkConn *splunk.SplunkConnection

func Connect2SplunkServer() (*splunk.SplunkConnection, error) {
	splunkHost, err := GetSplunkHostByEnvRole([]string{"prod", "MST"})
	if err != nil {
		return &splunk.SplunkConnection{}, err
	}
	var conf settings.Configuration
	conf = settings.GetConfig()
	splunkConn := splunk.SplunkConnection{
		Username:   conf.SplunkAdminUsername,
		Password:   conf.SplunkAdminPassword,
		BaseURL:    "https://" + splunkHost.Name,
		SplunkUser: "admin",
		SplunkApp:  "admin",
	}
	_, err = splunkConn.Login()
	if err != nil {
		return &splunk.SplunkConnection{}, err
	}

	fmt.Printf("%s Connected to splunk server. BaseURL:%v, Username:%v, Splunk App:%v ... [OK]\n", time.Now().Format("2006-01-02 15:04:05"), splunkConn.BaseURL, splunkConn.Username, splunkConn.SplunkApp)
	return &splunkConn, nil
}

func GetLogSize(month, serviceID string) ([]LogSize, error) {
	var err error
	if splunkConn == nil {
		splunkConn, err = Connect2SplunkServer()
		if err != nil {
			return nil, err
		}
	}

	statement := "search index=summary search_name=sum_license_usage2 earliest=-6mon@mon latest=now month=" +
		month + "* service_id=" + serviceID + " | table month service_id h virtual_mb | sort - virtual_mb"
	_, events, err := splunkConn.Search(statement, map[string]string{"maxout": "0"})
	if err != nil {
		return nil, err
	}
	var logSize LogSize
	var logSizes []LogSize
	host := map[string]bool{} // to eliminate duplicated result

	for _, event := range events {
		var f interface{}
		json.Unmarshal([]byte(event), &f)
		obj := f.(map[string]interface{})
		result := obj["result"].(map[string]interface{})
		logSize.Host = result["h"].(string)
		sizeStr := result["virtual_mb"].(string)
		logSize.ServiceID = serviceID
		value, _ := strconv.ParseFloat(sizeStr, 32)
		logSize.SizeMB = int(value)
		if host[logSize.Host] == false {
			logSizes = append(logSizes, logSize)
		}
		host[logSize.Host] = true
	}
	return logSizes, err
}

func GetStorageSize(month, serviceID string) ([]StorageSize, error) {
	var err error
	if splunkConn == nil {
		splunkConn, err = Connect2SplunkServer()
		if err != nil {
			return nil, err
		}
	}
	statement := "search index=summary search_name=sum_diskusage* earliest=-6mon@mon latest=now month=" +
		month + "* service_id=" + serviceID + " | table month service_id index_dirname virtual_mb | sort - virtual_mb"
	// settings.WriteDebugLog(statement)
	_, events, err := splunkConn.Search(statement, map[string]string{"maxout": "0"})
	if err != nil {
		return nil, err
	}
	var storage StorageSize
	var storages []StorageSize

	// settings.WriteDebugLog(events)
	index := map[string]bool{} // to eliminate duplicated result
	for _, event := range events {
		var f interface{}
		json.Unmarshal([]byte(event), &f)
		obj := f.(map[string]interface{})
		result := obj["result"].(map[string]interface{})
		storage.IndexName = result["index_dirname"].(string)
		sizeStr := result["virtual_mb"].(string)
		storage.ServiceID = serviceID
		value, _ := strconv.ParseFloat(sizeStr, 32)
		storage.SizeMB = int(value)
		if index[storage.IndexName] == false {
			storages = append(storages, storage)
		}
		index[storage.IndexName] = true
	}
	// settings.WriteDebugLog(storages)
	return storages, err
}
