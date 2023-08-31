package handlers

import (
	"net/http"
)

func setClientBadRequest(w http.ResponseWriter, err string) {
	clientErrMessage(w, http.StatusBadRequest, "domain", err)
}

type jsonReponse struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

func clientErrMessage(w http.ResponseWriter, statusCode int, errtype string, msg string) error {
	model := jsonReponse{
		Type:  errtype,
		Error: msg,
	}
	return setResponseJSON(w, statusCode, model)
}
