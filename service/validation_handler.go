package main

import (
	"encoding/json"
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
	result, msg, err := h.checkFunc(email)

	resultMap := map[string]interface{}{}
	resultMap["result"] = result
	resultMap["msg"] = msg
	if err != nil {
		resultMap["err"] = err.Error()
	}
	b, _ := json.MarshalIndent(resultMap, "", "  ")
	w.Write(b)
}
