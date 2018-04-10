package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"github.com/srinathgs/mysqlstore"
	"time"
)

var mysqldb *gorm.DB
var store *mysqlstore.MySQLStore

func Connect2Database() error {
	var conf settings.Configuration
	var err error
	conf = settings.GetConfig()

	mysqldb, err = gorm.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp(["+conf.DBHost+"]:"+conf.DBPort+")/"+conf.DBDatabase+"?charset=utf8&parseTime=True&loc=Local")
	// defer mysqldb.Close() //sql: database is closed
	if err != nil {
		return err
	}
	// mysqldb.DB().SetConnMaxLifetime(time.Minute * 5)
	mysqldb.DB().SetMaxIdleConns(0)
	// mysqldb.DB().SetMaxOpenConns(5)
	mysqldb.LogMode(conf.GormLogMode)
	fmt.Printf("%s Connected to database server. Host:%v, Port:%v, User:%v, Database:%v ... [OK]\n", time.Now().Format("2006-01-02 15:04:05"), conf.DBHost, conf.DBPort, conf.DBUser, conf.DBDatabase)

	return PrepareSessionStore(conf)
}

func PrepareSessionStore(conf settings.Configuration) error {
	var err error
	store, err = mysqlstore.NewMySQLStore(conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBHost+":"+conf.DBPort+")/"+conf.DBDatabase+"?charset=utf8&parseTime=True&loc=Local", "sessions", "/", conf.CookieTimeSecond, []byte(""))
	if err != nil {
		return err
	}
	fmt.Printf("%s %s", time.Now().Format("2006-01-02 15:04:05"), "Prepared session store ... [OK]\n")
	return nil
	//defer store.Close() //this will cause Error: sql: statement is closed
}
