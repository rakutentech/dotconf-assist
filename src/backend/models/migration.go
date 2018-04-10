package models

import (
	"github.com/jinzhu/gorm"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

// Migrate models
func Migrate(conffile string) error {
	if err := settings.ReadConfigFile(conffile); err != nil {
		log.Fatal(err.Error())
		return err
	}
	conf := settings.GetConfig()
	conn, err := gorm.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp(["+conf.DBHost+"]:"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	if res := conn.AutoMigrate(&User{}, &SplunkUser{}, &SplunkHost{}, &Announcement{}, &Forwarder{}, &ServerClass{}, &FileInput{}, &ScriptInput{}, &UnixAppInput{}, &App{}, &Deployment{}, &ServerClassForwarder{}, &UnitPrice{}); res.Error != nil {
		log.Fatal(res.Error.Error())
		return res.Error
	}

	// if res := conn.Model(&SplunkUser{}).AddForeignKey("user_name", "users(user_name)", "CASCADE", "CASCADE"); res.Error != nil {
	// 	err := res.Error
	// 	if strings.Index(err.Error(), "already exists for table") == -1 {
	// 		return err
	// 	}
	// }

	// if res := conn.Model(&Forwarder{}).AddForeignKey("user_name", "users(user_name)", "CASCADE", "CASCADE"); res.Error != nil {
	// 	err := res.Error
	// 	if strings.Index(err.Error(), "already exists for table") == -1 {
	// 		return err
	// 	}
	// }

	// if res := conn.Model(&Forwarder{}).AddUniqueIndex("uniq_idx_name_user", "name", "user_name"); res.Error != nil {
	// 	err := res.Error
	// 	if strings.Index(err.Error(), "already exists for table") == -1 {
	// 		return err
	// 	}
	// }

	// if res := conn.Model(&ServerClass{}).AddUniqueIndex("uniq_idx_name_user", "name", "user_name"); res.Error != nil {
	// 	err := res.Error
	// 	if strings.Index(err.Error(), "already exists for table") == -1 {
	// 		return err
	// 	}
	// }

	// if res := conn.Model(&ServerClassForwarder{}).AddUniqueIndex("uniq_idx_scid_fwdrid", "server_class_id", "forwarder_id"); res.Error != nil {
	// 	err := res.Error
	// 	if strings.Index(err.Error(), "already exists for table") == -1 {
	// 		return err
	// 	}
	// }

	// Save first admin
	pass, err := bcrypt.GenerateFromPassword([]byte(conf.AdminPassword), bcrypt.DefaultCost)
	admin := User{
		UserName: conf.AdminUsername,
		Email:    conf.AdminEmail,
		Admin:    true,
		Status:   "Approved",
		Password: string(pass),
	}

	if res := SaveUser(admin, false); res != nil {
		if strings.Index(res.Error(), "Duplicate key in container") == -1 {
			return res
		}
	}

	conn.Close()
	return nil
}
