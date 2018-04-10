package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"io"
	"io/ioutil"
	"net/http"
	// "reflect"
	"runtime"
	"strings"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SetResponseHeaders(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	// settings.WriteInfoLog(r.Header.Get("Origin"))
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	// w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "x-auth-token, Cookie, X-CSRFToken, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Cache-Control", "no-cache")
}

func WriteErrorResponse(w http.ResponseWriter, errCode int, err error) {
	w.WriteHeader(errCode)
	result := Result{Code: errCode, Msg: err.Error()}
	bytes, _ := json.Marshal(result)
	w.Write(bytes)
}

func WriteResponse(w http.ResponseWriter, code int, v interface{}) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.WriteHeader(code) //header
	if v != nil {
		w.Write(res) //body data
	} else {
		result := Result{Code: code, Msg: "OK"}
		bytes, _ := json.Marshal(result)
		w.Write(bytes)
	}
	return nil
}

func WriteBlobResponse(w http.ResponseWriter, code int, v []byte) error {
	w.WriteHeader(code) //header
	if v != nil {
		w.Write(v) //body data
	} else {
		result := Result{Code: code, Msg: "OK"}
		bytes, _ := json.Marshal(result)
		w.Write(bytes)
	}
	return nil
}

func ParseRequest(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10485760)) //10MB
	if err != nil {
		return []byte(""), err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body)) //add this line to restore the body which will be used in logger.go again

	err = r.Body.Close()
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}

func ValidateToken(tokenStr string) (string, string, *jwt.Token) {
	if tokenStr == "" {
		return "No token provided", "", nil
	}
	conf = settings.GetConfig()
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(conf.TokenSecret), nil
	})

	if token.Valid {
		userName, _ := token.Claims.(jwt.MapClaims)["sub"].(string)
		// email, _ := token.Claims["sub"].(string)
		return "Token valid", userName, token
	} else if ve, ok := err.(*jwt.ValidationError); ok { //invalid
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "Token invalid", "", token
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return "Token expired or not active yet", "", token
		} else {
			return "Cannot handle token", "", token
		}
	} else {
		return "Cannot handle token", "", token
	}
}

func GetCurrentFuncName() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	dotPos := strings.LastIndex(f.Name(), ".")
	return f.Name()[dotPos+1:]
}

func ConvertEmpty2Percent(str string) string {
	if str == "" {
		return "%"
	} else {
		return str
	}
}

func ConcertNull2EmptyStr(attr *string) {
	if *attr == "null" {
		*attr = ""
	}
}

func PreflightHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	settings.WriteLog(r, http.StatusOK, nil, GetCurrentFuncName())
}
