package models

import (
	"errors"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
)

func IsAdmin(userName string) bool {
	var user User
	res := mysqldb.Where("user_name = ? ", userName).Find(&user)
	if res.Error != nil { //record not found
		return false
	}
	return user.Admin
}

func ResetUserPassword(user User) error {
	usr, err := GetUser(user.UserName)
	if err != nil {
		return errors.New("Incorrect username or email")
	}
	if usr.Email == user.Email {
		usr.Password = user.Password
		res := mysqldb.Save(&usr)
		if res.Error != nil {
			return res.Error
		}
	} else {
		return errors.New("Incorrect username or email")
	}
	return nil
}

func SaveUser(user User, sendEmail bool) error {
	res := mysqldb.Save(&user)
	if res.Error != nil {
		return res.Error
	}

	if sendEmail {
		var conf settings.Configuration
		conf = settings.GetConfig()
		go func() {
			SendEmail(conf.EmailSender, conf.AdminEmail, "admin", "admin", "request-account",
				"[SPaaS] [Need Action] Request "+conf.AppName+" account", conf.AppUrl, user, "")
			SendEmail(conf.EmailSender, user.Email, "user", user.UserName, "request-account",
				"[SPaaS] Request "+conf.AppName+" account", conf.AppUrl, user, "")
		}()
	}
	return nil
}

func GetUsers() ([]User, error) {
	// settings.WriteDebugLog(mysqldb.DB().Stats())
	var users []User
	res := mysqldb.Order("id").Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	return users, nil
}

func GetUser(userName string) (User, error) {
	var user User
	res := mysqldb.Where("user_name = ? ", userName).Find(&user)
	if res.Error != nil { //record not found
		return User{}, res.Error
	}
	return user, nil
}

func UpdateUser(userName string, newUser User, action string) error {
	var user User
	var oldUser User
	res := mysqldb.Where("user_name = ? ", userName).Find(&user)
	if res.Error != nil { //record not found
		return res.Error
	}
	oldUser = user
	if action == "change_status" {
		user.Status = newUser.Status
	} else if action == "change_password" {
		user.Password = newUser.Password
	} else {
		user.Email = newUser.Email
		user.EmailForEmergency = newUser.EmailForEmergency
		user.GroupName = newUser.GroupName
		user.AppTeamName = newUser.AppTeamName
		user.ServiceID = newUser.ServiceID
	}

	res = mysqldb.Save(&user)
	if res.Error != nil {
		return res.Error
	}

	go func() {
		var conf settings.Configuration
		conf = settings.GetConfig()
		if action == "change_status" { // approve or cancel
			if newUser.Status == "Canceled" {
				err := SendEmail(conf.EmailSender, conf.AdminEmail, "admin", "admin", "cancel-account",
					"[SPaaS] Updated "+conf.AppName+" account information", conf.AppUrl, newUser, "")
				if err != nil {
					settings.WriteInfoLog(err)
				}
				err = SendEmail(conf.EmailSender, newUser.Email, "user", newUser.UserName, "cancel-account",
					"[SPaaS] Updated "+conf.AppName+" account information", conf.AppUrl, newUser, "")
				if err != nil {
					settings.WriteInfoLog(err)
				}
			} else if newUser.Status == "Approved" {
				err := SendEmail(conf.EmailSender, conf.AdminEmail, "admin", "admin", "approve-account",
					"[SPaaS] Updated "+conf.AppName+" account information", conf.AppUrl, newUser, "")
				if err != nil {
					settings.WriteInfoLog(err)
				}
				err = SendEmail(conf.EmailSender, newUser.Email, "user", newUser.UserName, "approve-account",
					"[SPaaS] Updated "+conf.AppName+" account information", conf.AppUrl, newUser, "")
				if err != nil {
					settings.WriteInfoLog(err)
				}
			}
		} else {
			err := SendEmail(conf.EmailSender, conf.AdminEmail, "common", "admin", "update-account",
				"[SPaaS] Updated "+conf.AppName+" account information", conf.AppUrl, newUser, oldUser)
			if err != nil {
				settings.WriteInfoLog(err)
			}
			err = SendEmail(conf.EmailSender, newUser.Email, "common", newUser.UserName, "update-account",
				"[SPaaS] Updated "+conf.AppName+" account information", conf.AppUrl, newUser, oldUser)
			if err != nil {
				settings.WriteInfoLog(err)
			}
		}
	}()
	return nil
}

func DeleteUser(userName string) error {
	user, err := GetUser(userName)
	if err != nil {
		return err
	}
	return mysqldb.Delete(&user).Error
}
