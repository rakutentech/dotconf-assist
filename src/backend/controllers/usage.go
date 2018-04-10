package controllers

import (
	"errors"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"net/http"
)

func GetUsagesHandler(w http.ResponseWriter, r *http.Request) { //get all usages
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var logSizes []models.LogSize
	var storageSizes []models.StorageSize
	usageType := r.URL.Query().Get("type")
	month := r.URL.Query().Get("month")
	serviceID := r.URL.Query().Get("serviceid")
	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if usageType == "log_size" {
		if logSizes, err = models.GetLogSize(month, serviceID); err != nil {
			goto Error
		}
		if err = WriteResponse(w, http.StatusOK, logSizes); err != nil {
			goto Error
		}
	} else {
		if storageSizes, err = models.GetStorageSize(month, serviceID); err != nil {
			goto Error
		}
		if err = WriteResponse(w, http.StatusOK, storageSizes); err != nil {
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
