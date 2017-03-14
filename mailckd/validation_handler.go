package main

import (
	"encoding/json"
	"fmt"
	"github.com/smancke/mailck"
	"github.com/tarent/lib-compose/logging"
	"net/http"
)

// MailValidationFunction checks the checkEmail
type MailValidationFunction func(checkEmail string) (result mailck.Result, err error)

// ValidationHandler is a REST handler for mail validation.
type ValidationHandler struct {
	checkFunc MailValidationFunction
}

func NewValidationHandler(checkFunc MailValidationFunction) *ValidationHandler {
	return &ValidationHandler{
		checkFunc: checkFunc,
	}
}

func (h *ValidationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()

	email := r.Form.Get("mail")
	if email == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"result": "error", "resultDetail": "clientError", "message": "missing parameter: mail"}`)
		return
	}

	result, err := h.checkFunc(email)

	if err != nil {
		logging.Application(r.Header).WithError(err).WithField("mail", email).Info("check error")
		if result == mailck.ServiceError {
			w.WriteHeader(500)
		}
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	w.Write(b)
}
