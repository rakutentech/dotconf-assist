package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Log struct {
	Time         string
	Message      string
	ResponseCode int
	Action       string
	Method       string
	URI          string
	User         string
	RequestData  string
}

func maskPassword(str string) string {
	re := regexp.MustCompile("\"password\":\\s*\"([^\"]+)\"")
	return re.ReplaceAllString(str, "\"password\": \"*\"")
}

func getUserEmail(token_str string) string {
	if token_str == "" {
		return ""
	}
	conf = GetConfig()
	token, _ := jwt.Parse(token_str, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(conf.TokenSecret), nil
	})

	if token.Valid {
		email, _ := token.Claims.(jwt.MapClaims)["sub"].(string)
		return email
	} else {
		return ""
	}
}

func WriteLog(r *http.Request, code int, erro error, action string) {
	if r.Method == "OPTIONS" {
		return
	}

	msg := "OK"
	if erro != nil {
		msg = erro.Error()
	}

	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	bodyString := string(bodyBytes)
	bodyString = maskPassword(bodyString)
	// if len(bodyString) > 100 {
	// 	bodyString = bodyString[0:99]
	// }

	token := r.Header.Get("x-auth-token")
	email := getUserEmail(token)

	logs := Log{
		time.Now().Format("2006-01-02 15:04:05"), msg, code, action, r.Method, r.RequestURI, email, bodyString,
	}

	b, err := json.Marshal(logs)
	if err != nil {
		log.Fatal(err)
	}

	conf = GetConfig()
	if conf.LogIndent {
		var out bytes.Buffer
		json.Indent(&out, b, "", " ")
		fmt.Printf("%s\n", out.String())
	} else {
		fmt.Printf("%s\n", string(b))
	}
}

func WriteInfoLog(v interface{}) {
	log.Printf("[Info] %v\n", v)
}

func WriteDebugLog(v interface{}) {
	log.Printf("[Debug] %v\n", v)
}
