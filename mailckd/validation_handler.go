package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/smancke/mailck"
	"github.com/tarent/lib-compose/logging"
	"io/ioutil"
	"net/http"
	"strings"
)

type parameters struct {
	Mail    string `json:"mail"`
	Timeout string `json:"timeout"`
}

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

	if r.Method != "GET" && r.Method != "POST" {
		writeError(w, 405, "clientError", "method not allowed")
		return
	}

	if r.Method == "POST" &&
		!(r.Header.Get("Content-Type") == "application/json" ||
			r.Header.Get("Content-Type") == "application/x-www-form-urlencoded") {
		writeError(w, 415, "clientError", "Unsupported Media Type")
		return
	}

	if !strings.HasSuffix(r.URL.Path, "/validate") {
		writeError(w, 404, "clientError", "ressource not found")
		return
	}

	p, err := h.readParameters(r)
	if err != nil {
		writeError(w, 400, "clientError", err.Error())
		return
	}

	result, err := h.checkFunc(p.Mail)

	if err != nil {
		logging.Application(r.Header).WithError(err).WithField("mail", p.Mail).Info("check error")
		if result == mailck.MailserverError {
			w.WriteHeader(502)
		} else {
			w.WriteHeader(500)
		}
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	w.Write(b)
}

func (h *ValidationHandler) readParameters(r *http.Request) (parameters, error) {
	p := parameters{}

	if r.Header.Get("Content-Type") == "application/json" {
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &p)
		if err != nil {
			return p, err
		}
	}

	// overwrite by form paramters, if any
	r.ParseForm()
	if r.Form.Get("mail") != "" {
		p.Mail = r.Form.Get("mail")
	}
	if r.Form.Get("timeout") != "" {
		p.Timeout = r.Form.Get("timeout")
	}

	if p.Mail == "" {
		return p, errors.New("missing parameter: mail")
	}

	return p, nil
}

func writeError(w http.ResponseWriter, code int, resultDetail, message string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"result": "error", "resultDetail": "%v", "message": "%v"}`, resultDetail, message)
}
