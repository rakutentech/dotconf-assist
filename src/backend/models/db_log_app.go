package models

import (
// "github.com/rakutentech/dotconf-assist/src/backend/settings"
)

func SaveApp(app App) error {
	// user, err := GetUser(userName)
	// if err != nil {
	// 	return err
	// }
	// app.UserID = user.ID
	// app.SplunkHostID = -1
	app.DeployStatus = -1 // not configured
	res := mysqldb.Save(&app)
	if res.Error != nil {
		return res.Error
	}
	return BindAppWithInputs(&app)
}

func BindAppWithInputs(app *App) error {
	var inputIDs []int
	if !app.UnixApp {
		for _, ids := range app.FileInputIDs {
			inputIDs = append(inputIDs, ids.ID)
		}
		UpdateInputFields(inputIDs, []string{"app_id"}, []interface{}{app.ID}, "file_inputs")
		inputIDs = inputIDs[:0]
		for _, ids := range app.ScriptInputIDs {
			inputIDs = append(inputIDs, ids.ID)
		}
		UpdateInputFields(inputIDs, []string{"app_id"}, []interface{}{app.ID}, "script_inputs")
	} else {
		for _, ids := range app.UnixAppInputIDs {
			inputIDs = append(inputIDs, ids.ID)
		}
		UpdateInputFields(inputIDs, []string{"app_id"}, []interface{}{app.ID}, "unix_app_inputs")
	}
	return nil
}

func unbindAppFromInputs(unixApp bool, appID int) error {
	if !unixApp {
		UpdateInputsFieldsByAppID(appID, []string{"app_id"}, []interface{}{-1}, "file_inputs")
		UpdateInputsFieldsByAppID(appID, []string{"app_id"}, []interface{}{-1}, "script_inputs")
	} else {
		UpdateInputsFieldsByAppID(appID, []string{"app_id"}, []interface{}{-1}, "unix_app_inputs")
	}
	return nil
}

func GetApps(envUser []string, getInputs bool, getServerClasses bool, getForwarders bool, isAdmin bool) ([]App, error) {
	var apps []App
	// user, err := GetUser(envUser[1])
	// if err != nil {
	// 	return nil, err
	// }
	// serverClass := map[int][]ServerClass{}

	if isAdmin {
		res := mysqldb.Where("env = ?", envUser[0]).Order("id").Find(&apps)
		if res.Error != nil { // res.Error is nil even if no record found
			return nil, res.Error
		}
	} else {
		// res := mysqldb.Where("env = ? AND user_id = ?", envUser[0], user.ID).Order("id").Find(&apps)
		res := mysqldb.Where("env = ? AND user_name = ?", envUser[0], envUser[1]).Order("id").Find(&apps)
		if res.Error != nil { // res.Error is nil even if no record found
			return nil, res.Error
		}
	}

	if getInputs {
		for i, app := range apps {
			if !app.UnixApp {
				inputs, err := GetInputsByAppID(envUser, app.ID, "file", isAdmin)
				if err != nil {
					return nil, err
				}
				apps[i].FileInputs = inputs.([]FileInput)
				inputs, err = GetInputsByAppID(envUser, app.ID, "script", isAdmin)
				if err != nil {
					return nil, err
				}
				apps[i].ScriptInputs = inputs.([]ScriptInput)
			} else {
				inputs, err := GetInputsByAppID(envUser, app.ID, "unixapp", isAdmin)
				if err != nil {
					return nil, err
				}
				apps[i].UnixAppInputs = inputs.([]UnixAppInput)
			}
		}
	}

	if getServerClasses {
		for i, app := range apps {
			scIDs, err := GetServerClassIDsByAppID(app.ID, getForwarders)
			if err != nil {
				return nil, err
			}
			apps[i].ServerClass, err = GetServerClassesByIDs(scIDs, true)
		}
	}
	return apps, nil
}

func GetApp(envUser []string, appID int, getInputs bool, isAdmin bool) (App, error) {
	var app App
	res := mysqldb.Where("id = ?", appID).Find(&app)
	if res.Error != nil { //record not found
		return App{}, res.Error
	}

	if getInputs {
		if app.UnixApp {
			inputs, err := GetInputsByAppID(envUser, app.ID, "file", isAdmin)
			if err != nil {
				return app, err
			}
			app.FileInputs = inputs.([]FileInput)

			inputs, err = GetInputsByAppID(envUser, app.ID, "script", isAdmin)
			if err != nil {
				return app, err
			}
			app.ScriptInputs = inputs.([]ScriptInput)
		} else {
			inputs, err := GetInputsByAppID(envUser, app.ID, "unixapp", isAdmin)
			if err != nil {
				return app, err
			}
			app.UnixAppInputs = inputs.([]UnixAppInput)
		}
	}
	return app, nil
}

func GetAppName(appID int) (string, error) {
	var app App
	res := mysqldb.Select("name").Where("id = ?", appID).Find(&app)
	if res.Error != nil { //record not found
		return "", res.Error
	}
	return app.Name, nil
}

func UpdateApp(appID int, newApp App) error {
	var app App
	res := mysqldb.Where("id = ?", appID).Find(&app)
	if res.Error != nil { //record not found
		return res.Error
	}
	oldAppName := app.Name
	app.Name = newApp.Name
	res = mysqldb.Save(&app)
	if res.Error != nil {
		return res.Error
	}
	if oldAppName == newApp.Name { //name not changed, update list
		newApp.ID = app.ID
		unbindAppFromInputs(app.UnixApp, appID)
		BindAppWithInputs(&newApp)
	}
	return nil
}

func UpdateAppFieldsByID(appID int, fields []string, values []interface{}, tableName string) error { //call when change app status
	keyValue := map[string]interface{}{}
	for i, _ := range fields {
		keyValue[fields[i]] = values[i]
	}
	res := mysqldb.Table(tableName).Where("id = ?", appID).Updates(keyValue)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteApp(appID int) error {
	app, err := GetApp([]string{"", ""}, appID, false, false)
	if err != nil {
		return err
	}
	unbindAppFromInputs(app.UnixApp, appID)
	return mysqldb.Delete(&app).Error
}
