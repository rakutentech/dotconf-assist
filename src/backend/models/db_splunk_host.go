package models

import (
// "github.com/rakutentech/dotconf-assist/src/backend/settings"
)

func SaveSplunkHost(splunkHost SplunkHost) error {
	res := mysqldb.Save(&splunkHost)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func GetSplunkHosts(envAndRole []string) ([]SplunkHost, error) {
	var splunkHosts []SplunkHost
	res := mysqldb.Where("env LIKE ? AND role LIKE ?", envAndRole[0], envAndRole[1]).Order("name desc").Find(&splunkHosts)
	if res.Error != nil {
		return nil, res.Error
	}

	return splunkHosts, nil
}

func GetSplunkHost(name string) (SplunkHost, error) {
	var splunkHost SplunkHost
	res := mysqldb.Where("name = ? ", name).Find(&splunkHost)
	if res.Error != nil { //record not found
		return SplunkHost{}, res.Error
	}
	return splunkHost, nil
}

func GetSplunkHostByEnvRole(envAndRole []string) (SplunkHost, error) {
	var splunkHost SplunkHost
	res := mysqldb.Where("env = ? AND role = ?", envAndRole[0], envAndRole[1]).Find(&splunkHost)
	if res.Error != nil { //record not found
		return SplunkHost{}, res.Error
	}
	return splunkHost, nil
}

func UpdateSplunkHost(name string, newSplunkHost SplunkHost) error {
	var splunkHost SplunkHost
	res := mysqldb.Where("name = ? ", name).Find(&splunkHost)
	if res.Error != nil { //record not found
		return res.Error
	}
	splunkHost.Name = newSplunkHost.Name
	splunkHost.Role = newSplunkHost.Role
	splunkHost.Env = newSplunkHost.Env
	return SaveSplunkHost(splunkHost)
}

func DeleteSplunkHost(name string) error {
	splunkHost, err := GetSplunkHost(name)
	if err != nil {
		return err
	}
	return mysqldb.Delete(&splunkHost).Error
}
