package main

import (
	"fmt"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/routers"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	conffile := "src/backend/settings/conf.json"
	if err := settings.ReadConfigFile(conffile); err != nil {
		log.Fatal(err.Error())
		return
	}
	if err := models.Connect2Database(); err != nil {
		log.Fatal(err.Error())
		return
	}

	for _, arg := range os.Args {
		if arg == "--with-migration" {
			models.Migrate(conffile)
		}
	}

	router := routers.NewRouter()

	conf := settings.GetConfig()
	if conf.EnableHTTPS {
		fmt.Printf("%s Enabling HTTPS ... [OK]\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("%s Listening on port 4201 ... [OK]\n", time.Now().Format("2006-01-02 15:04:05"))
		log.Fatal(http.ListenAndServeTLS(":4201", conf.Cert, conf.Key, router))
	} else {
		fmt.Printf("%s Listening on port 4201 ... [OK]\n", time.Now().Format("2006-01-02 15:04:05"))
		log.Fatal(http.ListenAndServe(":4201", router))
	}
}
