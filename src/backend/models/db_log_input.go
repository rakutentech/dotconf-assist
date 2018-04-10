package models

import (
// "errors"
// "github.com/jinzhu/gorm"
// "github.com/rakutentech/dotconf-assist/src/backend/settings"
)

func SaveFileInput(input FileInput, envUser []string) error {
	// splunkUser, err := GetSplunkUser(envUser[1], envUser[0])
	// if err != nil { //Splunk user not created yet
	// 	return errors.New("You have not created Splunk user yet")
	// }
	// _ = mysqldb.Exec("INSERT INTO inputs(script) VALUES(?)", input.Script)
	// input.SplunkUserID = splunkUser.ID
	input.AppID = -1
	res := mysqldb.Save(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func SaveScriptInput(input ScriptInput, envUser []string) error {
	// splunkUser, err := GetSplunkUser(envUser[1], envUser[0])
	// if err != nil { //Splunk user not created yet
	// 	return errors.New("You have not created Splunk user yet")
	// }
	// _ = mysqldb.Exec("INSERT INTO inputs(script) VALUES(?)", input.Script)
	// input.SplunkUserID = splunkUser.ID
	input.AppID = -1
	res := mysqldb.Save(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func SaveUnixAppInput(inputs []UnixAppInput, envUser []string) error {
	// splunkUser, err := GetSplunkUser(envUser[1], envUser[0])
	// if err != nil { //Splunk user not created yet
	// 	return errors.New("You have not created Splunk user yet")
	// }
	// _ = mysqldb.Exec("INSERT INTO inputs(script) VALUES(?)", input.Script)
	for _, input := range inputs {
		// input.SplunkUserID = splunkUser.ID
		input.AppID = -1
		res := mysqldb.Save(&input)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func GetInputs(envUser []string, inputType string, isAdmin bool) (interface{}, error) {
	// var err error
	// var splunkUser SplunkUser
	// if !isAdmin {
	// 	splunkUser, err = GetSplunkUser(envUser[1], envUser[0])
	// 	if err != nil {
	// 		return nil, errors.New("You have not created Splunk user yet")
	// 	}
	// }

	var fileInputs []FileInput
	var scriptInputs []ScriptInput
	var unixAppInputs []UnixAppInput
	appName := map[int]string{}

	switch inputType {
	case "file":
		{
			if isAdmin {
				res := mysqldb.Order("id").Find(&fileInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			} else {
				// res := mysqldb.Where("splunk_user_id = ?", splunkUser.ID).Order("id").Find(&fileInputs)
				res := mysqldb.Where("env = ? AND user_name = ?", envUser[0], envUser[1]).Order("id").Find(&fileInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			}

			for i, input := range fileInputs {
				if input.AppID != -1 {
					if appName[input.AppID] == "" {
						fileInputs[i].AppName, _ = GetAppName(input.AppID)
						appName[input.AppID] = fileInputs[i].AppName
					} else {
						fileInputs[i].AppName = appName[input.AppID]
					}
				} else {
					fileInputs[i].AppName = ""
				}
			}
			return fileInputs, nil
		}
	case "script":
		{
			if isAdmin {
				res := mysqldb.Order("id").Find(&scriptInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			} else {
				// res := mysqldb.Where("splunk_user_id = ?", splunkUser.ID).Order("id").Find(&scriptInputs)
				res := mysqldb.Where("env = ? AND user_name = ?", envUser[0], envUser[1]).Order("id").Find(&scriptInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			}

			for i, input := range scriptInputs {
				if input.AppID != -1 {
					if appName[input.AppID] == "" {
						scriptInputs[i].AppName, _ = GetAppName(input.AppID)
						appName[input.AppID] = scriptInputs[i].AppName
					} else {
						scriptInputs[i].AppName = appName[input.AppID]
					}
				} else {
					scriptInputs[i].AppName = ""
				}
			}
			return scriptInputs, nil
		}
	case "unixapp":
		{
			if isAdmin {
				res := mysqldb.Order("id").Find(&unixAppInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			} else {
				// res := mysqldb.Where("splunk_user_id = ?", splunkUser.ID).Order("id").Find(&unixAppInputs)
				res := mysqldb.Where("env = ? AND user_name = ?", envUser[0], envUser[1]).Order("id").Find(&unixAppInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			}

			for i, input := range unixAppInputs {
				if input.AppID != -1 {
					if appName[input.AppID] == "" {
						unixAppInputs[i].AppName, _ = GetAppName(input.AppID)
						appName[input.AppID] = unixAppInputs[i].AppName
					} else {
						unixAppInputs[i].AppName = appName[input.AppID]
					}
				} else {
					unixAppInputs[i].AppName = ""
				}
			}
			return unixAppInputs, nil
		}
	default:
		{
		}
	}
	return nil, nil
}

func GetInputsByAppID(envUser []string, appID int, inputType string, isAdmin bool) (interface{}, error) { //used in app view
	// var err error
	// var splunkUser SplunkUser
	// if !isAdmin {
	// 	splunkUser, err = GetSplunkUser(envUser[1], envUser[0])
	// 	if err != nil {
	// 		return nil, errors.New("You have not created Splunk user yet")
	// 	}
	// }

	var fileInputs []FileInput
	var scriptInputs []ScriptInput
	var unixAppInputs []UnixAppInput
	switch inputType {
	case "file":
		{
			if isAdmin {
				res := mysqldb.Select("id, log_file_path, sourcetype, log_file_size, data_retention_period, memo, blacklist, crcsalt").
					Where("app_id = ?", appID).Find(&fileInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			} else {
				res := mysqldb.Select("id, log_file_path, sourcetype, log_file_size, data_retention_period, memo, blacklist, crcsalt").
					// Where("app_id = ? AND splunk_user_id = ?", appID, splunkUser.ID).Find(&fileInputs)
					Where("app_id = ? AND env = ? AND user_name= ?", appID, envUser[0], envUser[1]).Find(&fileInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			}

			return fileInputs, nil
		}
	case "script":
		{
			if isAdmin {
				res := mysqldb.Select("id, sourcetype, log_file_size, data_retention_period, os, `interval`, script_name, `option`, exefile").
					Where("app_id = ?", appID).Find(&scriptInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			} else {
				res := mysqldb.Select("id, sourcetype, log_file_size, data_retention_period, os, `interval`, script_name, `option`, exefile").
					// Where("app_id = ? AND splunk_user_id = ?", appID, splunkUser.ID).Find(&scriptInputs)
					Where("app_id = ? AND env = ? AND user_name= ?", appID, envUser[0], envUser[1]).Find(&scriptInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			}
			return scriptInputs, nil
		}
	case "unixapp":
		{
			if isAdmin {
				res := mysqldb.Select("id, script_name, data_retention_period, `interval`").
					Where("app_id = ?", appID).Find(&unixAppInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			} else {
				res := mysqldb.Select("id, script_name, data_retention_period, `interval`").
					// Where("app_id = ? AND splunk_user_id = ?", appID, splunkUser.ID).Find(&unixAppInputs)
					Where("app_id = ? AND env = ? AND user_name= ?", appID, envUser[0], envUser[1]).Find(&unixAppInputs)
				if res.Error != nil {
					return nil, res.Error
				}
			}

			return unixAppInputs, nil
		}
	default:
		{
		}
	}
	return nil, nil
}

func GetInput(inputID int, inputType string) (interface{}, error) {
	var fileInput FileInput
	var scriptInput ScriptInput
	var unixAppInput UnixAppInput
	switch inputType {
	case "file":
		{
			res := mysqldb.Where("id = ?", inputID).Find(&fileInput)
			if res.Error != nil { //record not found
				return FileInput{}, res.Error
			}
			return fileInput, nil
		}
	case "script":
		{
			res := mysqldb.Where("id = ?", inputID).Find(&scriptInput)
			if res.Error != nil { //record not found
				return ScriptInput{}, res.Error
			}
			return scriptInput, nil
		}
	case "unixapp":
		{
			res := mysqldb.Where("id = ?", inputID).Find(&unixAppInput)
			if res.Error != nil { //record not found
				return UnixAppInput{}, res.Error
			}
			return unixAppInput, nil
		}
	}
	return FileInput{}, nil
}

func UpdateInputsFieldsByAppID(appID int, fields []string, values []interface{}, tableName string) error { //call when deleting app
	keyValue := map[string]interface{}{}
	for i, _ := range fields {
		keyValue[fields[i]] = values[i]
	}
	res := mysqldb.Table(tableName).Where("app_id = ?", appID).Updates(keyValue)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func UpdateInputFields(inputIDs []int, fields []string, values []interface{}, tableName string) error { //call when add app
	keyValue := map[string]interface{}{}
	for i, _ := range fields {
		keyValue[fields[i]] = values[i]
	}
	res := mysqldb.Table(tableName).Where("id IN (?)", inputIDs).Updates(keyValue)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func UpdateFileInput(inputID int, newInput FileInput) error {
	var input FileInput
	res := mysqldb.Where("id = ?", inputID).Find(&input)
	if res.Error != nil { //record not found
		return res.Error
	}
	input.LogFilePath = newInput.LogFilePath
	input.LogFileSize = newInput.LogFileSize
	input.Sourcetype = newInput.Sourcetype
	input.DataRetentionPeriod = newInput.DataRetentionPeriod
	input.AppID = newInput.AppID
	input.Blacklist = newInput.Blacklist
	input.Crcsalt = newInput.Crcsalt
	// input.Status = newInput.Status
	input.Memo = newInput.Memo
	res = mysqldb.Save(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func UpdateScriptInput(inputID int, newInput ScriptInput) error {
	var input ScriptInput
	res := mysqldb.Where("id = ?", inputID).Find(&input)
	if res.Error != nil { //record not found
		return res.Error
	}
	input.LogFileSize = newInput.LogFileSize
	input.Sourcetype = newInput.Sourcetype
	input.DataRetentionPeriod = newInput.DataRetentionPeriod
	input.OS = newInput.OS
	input.Interval = newInput.Interval
	input.ScriptName = newInput.ScriptName
	input.Option = newInput.Option
	input.Exefile = newInput.Exefile
	// input.AppID = newInput.AppID
	// input.Status = newInput.Status
	if len(newInput.Script) > 1 {
		input.Script = newInput.Script
	}
	res = mysqldb.Save(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func UpdateUnixAppInput(inputID int, newInput UnixAppInput) error {
	var input UnixAppInput
	res := mysqldb.Where("id = ?", inputID).Find(&input)
	if res.Error != nil { //record not found
		return res.Error
	}
	input.DataRetentionPeriod = newInput.DataRetentionPeriod
	input.Interval = newInput.Interval
	// input.Status = newInput.Status
	// input.AppID = newInput.AppID
	res = mysqldb.Save(&input)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteInput(inputID int, inputType string) error {
	input, err := GetInput(inputID, inputType)
	if err != nil { //not found
		return err
	}
	return mysqldb.Delete(input).Error
}
