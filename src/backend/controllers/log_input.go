package controllers

import (
	"encoding/json"
	"errors"
	// "fmt"
	"github.com/gorilla/mux"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	// "io"
	"bytes"
	"net/http"
	// "os"
	"mime/multipart"
	"strconv"
	// "time"
)

func AddInputHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var fileInput models.FileInput          //here is single object, not slice
	var scriptInput models.ScriptInput      //here is single object, not slice
	var unixAppInputs []models.UnixAppInput //slice
	env := r.URL.Query().Get("env")
	user := r.URL.Query().Get("user")
	inputType := r.URL.Query().Get("type")
	params := []string{ConvertEmpty2Percent(env), ConvertEmpty2Percent(user)}
	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	buf := new(bytes.Buffer) //Buffer is a struct
	var file multipart.File

	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if inputType == "file" {
		if body, err = ParseRequest(w, r); err != nil {
			goto Error
		}
		if err = json.Unmarshal(body, &fileInput); err != nil {
			goto Error
		}
		// fileInput.CreatedAt = time.Now()
		ConcertNull2EmptyStr(&fileInput.Blacklist)
		ConcertNull2EmptyStr(&fileInput.Crcsalt)
		ConcertNull2EmptyStr(&fileInput.Memo)
		if err = models.SaveFileInput(fileInput, params); err != nil {
			goto Error
		}
	} else if inputType == "script" {
		scriptInput.Sourcetype = r.FormValue("sourcetype")
		scriptInput.LogFileSize = r.FormValue("log_file_size")
		scriptInput.DataRetentionPeriod = r.FormValue("data_retention_period")
		scriptInput.OS = r.FormValue("os")
		scriptInput.Interval = r.FormValue("interval")
		var exeFile = r.FormValue("exefile")
		if exeFile == "true" {
			scriptInput.Exefile = true
		} else {
			scriptInput.Exefile = false
		}
		scriptInput.ScriptName = r.FormValue("script_name")
		scriptInput.Option = r.FormValue("option")
		scriptInput.Env = r.FormValue("env")
		scriptInput.UserName = r.FormValue("user_name")
		ConcertNull2EmptyStr(&scriptInput.Sourcetype)
		ConcertNull2EmptyStr(&scriptInput.ScriptName)
		ConcertNull2EmptyStr(&scriptInput.Option)

		file, _, err = r.FormFile("script")
		defer file.Close()
		if err != nil {
			responseCode = http.StatusBadRequest
			goto Error
		}
		buf.ReadFrom(file)
		scriptInput.Script = buf.Bytes() //convert Buffer to []byte
		// bufStr := buf.String()
		// fmt.Println(input.Script)
		// fmt.Println("%s\n", bufStr)
		// fmt.Println("%s\n", string(input.Script)) //ok

		if err = models.SaveScriptInput(scriptInput, params); err != nil {
			goto Error
		}
	} else if inputType == "unixapp" {
		if body, err = ParseRequest(w, r); err != nil {
			goto Error
		}
		if err = json.Unmarshal(body, &unixAppInputs); err != nil {
			goto Error
		}
		// fileInput.CreatedAt = time.Now()
		if err = models.SaveUnixAppInput(unixAppInputs, params); err != nil {
			goto Error
		}
	}

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetInputsHandler(w http.ResponseWriter, r *http.Request) { //get all inputs
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var idx int = 0
	var appID int = -1
	var anyInputs interface{}
	env := r.URL.Query().Get("env")
	user := r.URL.Query().Get("user")
	inputType := r.URL.Query().Get("type")
	appIDStr := r.URL.Query().Get("app_id")
	params := []string{ConvertEmpty2Percent(env), ConvertEmpty2Percent(user)}
	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	isAdmin := false

	if appIDStr != "" {
		appID, err = strconv.Atoi(appIDStr)
		if err != nil {
			goto Error
		}
	}

	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if userName == "admin" {
		isAdmin = true
	}

	if appIDStr == "" { //get all inputs, called in inputs view
		if anyInputs, err = models.GetInputs(params, inputType, isAdmin); err != nil {
			goto Error
		}
	} else { //get inputs with specific ID, bound(non -1) or not bound(-1), called in app's view
		if anyInputs, err = models.GetInputsByAppID(params, appID, inputType, isAdmin); err != nil {
			goto Error
		}
	}

	if inputType == "file" {
		if err = WriteResponse(w, http.StatusOK, anyInputs.([]models.FileInput)); err != nil {
			goto Error
		}
	} else if inputType == "script" {
		var scriptInputs []models.ScriptInput
		scriptInputs = anyInputs.([]models.ScriptInput)
		for idx, _ = range scriptInputs {
			scriptInputs[idx].ScriptCode = (string)(scriptInputs[idx].Script)
		}
		if err = WriteResponse(w, http.StatusOK, scriptInputs); err != nil {
			goto Error
		}
	} else if inputType == "unixapp" {
		if err = WriteResponse(w, http.StatusOK, anyInputs.([]models.UnixAppInput)); err != nil {
			goto Error
		}
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetInputHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var anyInput interface{}
	vars := mux.Vars(r)
	inputID, err := strconv.Atoi(vars["id"])
	inputType := r.URL.Query().Get("type")

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if anyInput, err = models.GetInput(inputID, inputType); err != nil {
		goto Error
	}

	if inputType == "file" {
		if err = WriteResponse(w, http.StatusOK, anyInput.(models.FileInput)); err != nil {
			goto Error
		}
	} else if inputType == "script" {
		if err = WriteResponse(w, http.StatusOK, anyInput.(models.ScriptInput)); err != nil {
			goto Error
		}
	} else if inputType == "unixapp" {
		if err = WriteResponse(w, http.StatusOK, anyInput.(models.UnixAppInput)); err != nil {
			goto Error
		}
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateInputHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var fileInput models.FileInput
	var scriptInput models.ScriptInput
	var unixAppInput models.UnixAppInput

	vars := mux.Vars(r)
	inputID, err := strconv.Atoi(vars["id"])
	inputType := r.URL.Query().Get("type")

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	buf := new(bytes.Buffer) //Buffer is a struct
	var file multipart.File

	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if inputType == "file" {
		if body, err = ParseRequest(w, r); err != nil {
			goto Error
		}
		if err = json.Unmarshal(body, &fileInput); err != nil {
			goto Error
		}
		ConcertNull2EmptyStr(&fileInput.Blacklist)
		ConcertNull2EmptyStr(&fileInput.Crcsalt)
		ConcertNull2EmptyStr(&fileInput.Memo)
		if err = models.UpdateFileInput(inputID, fileInput); err != nil {
			goto Error
		}
	} else if inputType == "script" {
		scriptInput.Sourcetype = r.FormValue("sourcetype")
		scriptInput.LogFileSize = r.FormValue("log_file_size")
		scriptInput.DataRetentionPeriod = r.FormValue("data_retention_period")
		scriptInput.OS = r.FormValue("os")
		scriptInput.Interval = r.FormValue("interval")
		var exeFile = r.FormValue("exefile")
		if exeFile == "true" {
			scriptInput.Exefile = true
		} else {
			scriptInput.Exefile = false
		}
		scriptInput.ScriptName = r.FormValue("script_name")
		scriptInput.Option = r.FormValue("option")
		ConcertNull2EmptyStr(&scriptInput.Option)

		file, _, err = r.FormFile("script")
		if err == nil {
			defer file.Close()
			buf.ReadFrom(file)
			scriptInput.Script = buf.Bytes() //convert Buffer to []byte
			// responseCode = http.StatusBadRequest
			// goto Error
		}
		if err = models.UpdateScriptInput(inputID, scriptInput); err != nil {
			goto Error
		}
	} else if inputType == "unixapp" {
		if body, err = ParseRequest(w, r); err != nil {
			goto Error
		}
		if err = json.Unmarshal(body, &unixAppInput); err != nil {
			goto Error
		}

		if err = models.UpdateUnixAppInput(inputID, unixAppInput); err != nil {
			goto Error
		}
	}

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func DeleteInputHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	vars := mux.Vars(r)
	inputID, err := strconv.Atoi(vars["id"])
	inputType := r.URL.Query().Get("type")

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if err = models.DeleteInput(inputID, inputType); err != nil {
		goto Error
	}

	// if inputType == "file" {
	// 	if err = models.DeleteFileInput(inputID); err != nil {
	// 		goto Error
	// 	}
	// } else if inputType == "script" {
	// 	if err = models.DeleteScriptInput(inputID); err != nil {
	// 		goto Error
	// 	}
	// } else if inputType == "unixapp" {
	// 	if err = models.DeleteUnixAppInput(inputID); err != nil {
	// 		goto Error
	// 	}
	// }

	responseCode = http.StatusOK
	WriteResponse(w, responseCode, nil)
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}
