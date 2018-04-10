package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rakutentech/dotconf-assist/src/backend/models"
	"github.com/rakutentech/dotconf-assist/src/backend/settings"
	"net/http"
	"time"
)

func AddAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var body []byte
	var responseCode int = http.StatusBadRequest
	var announcement models.Announcement //here is single object, not slice

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if !models.IsAdmin(userName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}

	if err = json.Unmarshal(body, &announcement); err != nil {
		goto Error
	}

	announcement.CreatedAt = time.Now()
	if err = models.SaveAnnouncement(announcement); err != nil {
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

func GetAnnouncementsHandler(w http.ResponseWriter, r *http.Request) { //get all announcements
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var announcements []models.Announcement

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}

	if announcements, err = models.GetAnnouncements(); err != nil {
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, announcements); err != nil {
		goto Error
	}
	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func GetAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var announcement models.Announcement

	vars := mux.Vars(r)
	announcementID := vars["id"]

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(userName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if announcement, err = models.GetAnnouncement(announcementID); err != nil {
		if err.Error() == "record not found" {
			responseCode = http.StatusNotFound
		}
		goto Error
	}

	if err = WriteResponse(w, http.StatusOK, announcement); err != nil {
		goto Error
	}

	responseCode = http.StatusOK
	settings.WriteLog(r, responseCode, nil, GetCurrentFuncName())
	return
Error:
	WriteErrorResponse(w, responseCode, err)
	settings.WriteLog(r, responseCode, err, GetCurrentFuncName())
}

func UpdateAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest
	var announcement models.Announcement //here is single object, not slice
	var body []byte

	vars := mux.Vars(r)
	announcementID := vars["id"]

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(userName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if body, err = ParseRequest(w, r); err != nil {
		goto Error
	}
	if err = json.Unmarshal(body, &announcement); err != nil {
		goto Error
	}

	if err = models.UpdateAnnouncement(announcementID, announcement); err != nil {
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

func DeleteAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeaders(w, r)
	var err error
	var responseCode int = http.StatusBadRequest

	vars := mux.Vars(r)
	announcementID := vars["id"]

	token := r.Header.Get("x-auth-token")
	msg, userName, _ := ValidateToken(token)
	if userName == "" {
		err = errors.New(msg)
		responseCode = http.StatusUnauthorized
		goto Error
	}
	if !models.IsAdmin(userName) {
		err = errors.New("Permission denied")
		responseCode = http.StatusForbidden
		goto Error
	}

	if err = models.DeleteAnnouncement(announcementID); err != nil {
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
