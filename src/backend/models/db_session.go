package models

import (
	"errors"
	"github.com/srinathgs/mysqlstore" //for session store
	"golang.org/x/crypto/bcrypt"
)

func GetSessionStore() (*mysqlstore.MySQLStore, error) {
	return store, nil
}

func VerifyCredentials(userName, password string) error {
	DeleteExpiredSession()
	var user User
	var admin User
	var err error
	res := mysqldb.Where("user_name = ?", userName).Find(&user)
	if res.Error != nil { //record not found
		return res.Error
	}

	if user.Status != "Approved" {
		return errors.New("Your account has not been approved yet")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		res := mysqldb.Where("user_name = ?", "admin").Find(&admin)
		if res.Error != nil { //record not found
			return err
		}
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteExpiredSession() {
	mysqldb.Exec("DELETE FROM sessions where expires_on < DATE_SUB(now(), interval 3 minute);")
}

// func heartBeat(){
// 	mysqldb.Exec("select 1;")
// }
