package main

import (
	"encoding/json"
	"fmt"
	"github.com/smancke/mailck"
	"net/http"
)

// MailValidationFunction checks the checkEmail
type MailValidationFunction func(checkEmail string) (result mailck.CheckResult, textMessage string, err error)

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

	email := r.Form.Get("email")
	if email == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"success": false, "msg": "missing parameter: email"}`)
		return
	}

	result, msg, err := h.checkFunc(email)

	resultMap := map[string]interface{}{
		"success": err == nil,
		"result":  result,
		"msg":     msg,
	}
	if err != nil {
		w.WriteHeader(500)
	}
	b, _ := json.MarshalIndent(resultMap, "", "  ")
	w.Write(b)
}
