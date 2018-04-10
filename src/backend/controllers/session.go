package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"github.com/srinathgs/mysqlstore"
	"net/http"
	"os"
	"strings"
	"time"
)

var conf settings.Configuration

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var body []byte
	var user models.User //user info in post form email + password(non-LDAP auth), or username + passowrd (LDAP)
	var usr models.User
	var tokenString string = ""
	var authOK bool = false

	if err = models.Connect2Database(); err != nil { //in case timeout
		goto Error
	}

	conf = settings.GetConfig()
	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}

	if err = json.Unmarshal(body, &user); err != nil {
		goto Error
	}
	authOK, err = authenticate(w, r, &user)
	if authOK {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.UserName,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * conf.TokenTimeHour).Unix(),
		})
		tokenString, err = token.SignedString([]byte(conf.TokenSecret))
		if err != nil {
			goto Error
		}

		usr, err = models.GetUser(user.UserName) //use usr to update,if not exist,to insert
		if err != nil {                          //login with LDAP and user is not in database
			usr.Admin = false
			usr.UserName = user.UserName
		}

		usr.LastLoginAt = time.Now()
		if err = models.SaveUser(usr, false); err != nil { //save user in database when user logged in
			goto Error
		}
		if err = setSession(&usr, w, r); err != nil {
			goto Error
		}
	} else {
		if !strings.Contains(err.Error(), "has not been approved") {
			err = errors.New("Authentication failed")
		}
		responseCode = http.StatusUnauthorized
		goto Error
	}

	w.Header().Set("x-auth-token", tokenString)
	if usr.Admin {
		w.Header().Set("Admin", "true")
	}
	w.Header().Set("Access-Control-Expose-Headers", "x-auth-token, Admin")

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var store *mysqlstore.MySQLStore
	var session *sessions.Session
	var loggedIn bool

	loggedIn, err = isLoggedin(w, r)
	if err != nil {
		goto Error
	}

	if !loggedIn {
		err = errors.New("Not logged in")
		goto Error
	}

	store, err = models.GetSessionStore()
	if err != nil {
		responseCode = http.StatusInternalServerError
		goto Error
	}
	session, err = getSession(w, r)
	if err != nil {
		goto Error
	}

	if err = store.Delete(r, w, session); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func ISO8601(t time.Time) string {
	var tz string
	zone, offset := t.Zone()
	if zone == "UTC" {
		tz = "Z"
	} else {
		if offset > 0 {
			tz = fmt.Sprintf("+%02d:00", offset/3600)
		} else {
			tz = fmt.Sprintf("%03d:00", offset/3600)
		}
	}
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d%s", year, month, day, hour, min, sec, tz)
}

func formatToken(exp, iss time.Time, method, email string) interface{} {
	expires_at := ISO8601(exp)
	issued_at := ISO8601(iss)
	b := []byte(`{"methods": ["`)
	b = append(b, method...)
	b = append(b, `"],"expires_at": "`...)
	b = append(b, expires_at...)
	b = append(b, `","extras": {},"user": {"domain": {"id": "default","name": "Default"},"id": "`...)
	b = append(b, email...)
	b = append(b, `","name": "`...)
	b = append(b, email...)
	b = append(b, `"},"audit_ids": ["ZzZwkUflQfygX7pdYDBCQQ"],"issued_at": "`...)
	b = append(b, issued_at...)
	b = append(b, `"}`...)

	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil
	}
	return f
}

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var token *jwt.Token
	var exp float64
	var iat float64
	var tokenStr = r.Header.Get("x-auth-token")
	type tokenObj struct {
		Token interface{} `json:"token"`
	}
	msg, email, token := ValidateToken(tokenStr)
	if email == "" {
		err = errors.New(msg)
		goto Error
	}

	exp, _ = token.Claims.(jwt.MapClaims)["exp"].(float64)
	iat, _ = token.Claims.(jwt.MapClaims)["iat"].(float64)

	if err = WriteResponse(w, http.StatusOK, tokenObj{formatToken(time.Unix(int64(exp), 0), time.Unix(int64(iat), 0), "token", email)}); err != nil {
		goto Error
	}
	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateSessionHandler(w http.ResponseWriter, r *http.Request) { //for frontend token update
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var tokenString string = ""
	var token *jwt.Token

	var tokenStr = r.Header.Get("x-auth-token")
	msg, email, _ := ValidateToken(tokenStr)
	if email == "" {
		err = errors.New(msg)
		goto Error
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * conf.TokenTimeHour).Unix(),
	})
	tokenString, err = token.SignedString([]byte(conf.TokenSecret))
	if err != nil {
		goto Error
	}

	w.Header().Set("x-auth-token", tokenString)
	if models.IsAdmin(email) {
		w.Header().Set("Admin", "true")
	}
	w.Header().Set("Access-Control-Expose-Headers", "x-auth-token, Admin")

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func isLoggedin(w http.ResponseWriter, r *http.Request) (loggedin bool, err error) {
	session, err := getSession(w, r)
	if err != nil {
		return false, err
	}

	if session.IsNew {
		return false, nil
	} else {
		loggedin = !session.IsNew
		// email = session.Values["email"].(string)
		return loggedin, nil
	}
}

func getOriginHostName(str string) string {
	hostAndPort := strings.Split(str, "//")
	if len(hostAndPort) < 2 {
		return ""
	}
	host := strings.Split(hostAndPort[1], ":")
	if len(host) < 1 {
		return hostAndPort[1]
	}
	return host[0]
}

func setSession(user *models.User, w http.ResponseWriter, r *http.Request) error {
	var err error
	var store *mysqlstore.MySQLStore
	store, err = models.GetSessionStore()
	if err != nil {
		return err
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
	}

	name, err := os.Hostname()
	if err != nil {
		return err
	}
	session, err := store.Get(r, getOriginHostName(origin)+"_"+name+"_session")
	session.Values["email"] = user.Email
	err = store.Save(r, w, session) //Set-Cookie session=... //or err = session.Save(r, w)
	return err
}

func getSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	store, err := models.GetSessionStore()
	if err != nil {
		return nil, err
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		// settings.WriteInfoLog("Request origin is empty")
	}
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	session, err := store.Get(r, getOriginHostName(origin)+"_"+name+"_session")
	return session, err
}

func authenticate(w http.ResponseWriter, r *http.Request, user *models.User) (bool, error) {
	conf = settings.GetConfig()
	err := models.VerifyCredentials(user.UserName, user.Password)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
