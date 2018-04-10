package settings

import (
	"encoding/json"
	"os"
	"time"
)

type Configuration struct {
	AppName             string
	AppUrl              string
	DBHost              string
	DBDatabase          string
	DBUser              string
	DBPassword          string
	DBPort              string
	TokenSecret         string
	CookieTimeSecond    int
	TokenTimeHour       time.Duration
	GormLogMode         bool
	TestMode            bool
	LogIndent           bool
	AdminUsername       string
	AdminPassword       string
	AdminEmail          string
	EmailSender         string
	SplunkAdminUsername string
	SplunkAdminPassword string
	DPLYSSHUser         string
	DPLYSSHPassword     string
	MailServer          string
	ProdSHC             string
	Proxy               string
	EnableHTTPS         bool
	Cert                string
	Key                 string
}

var conf Configuration

func ReadConfigFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	return err
}

func GetConfig() Configuration {
	return conf
}
