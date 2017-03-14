package main

import (
	"encoding/json"
	"errors"
	"github.com/smancke/mailck"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testValidationFunction(result mailck.Result, err error) MailValidationFunction {
	return func(checkEmail string) (mailck.Result, error) {
		if checkEmail != "foo@example.com" {
			panic("wrong email: " + checkEmail)
		}
		return result, err
	}
}

func Test_BadRequest(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.Valid, nil))

	req, err := http.NewRequest("POST", "/", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, "error", result["result"])
	assert.Equal(t, "clientError", result["resultDetail"])
	assert.Equal(t, "missing parameter: mail", result["message"])
}

func Test_SuccessPost(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.Valid, nil))
	req, err := http.NewRequest("POST", "http://localhost:3000/api/validate", strings.NewReader(`mail=foo@example.com`))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, "valid", result["result"])
	assert.Equal(t, "mailboxChecked", result["resultDetail"])
	assert.Equal(t, mailck.Valid.Message, result["message"])
}

func Test_SuccessGet(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.Valid, nil))
	req, err := http.NewRequest("GET", "http://localhost:3000/api/validate?mail=foo@example.com", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, "valid", result["result"])
	assert.Equal(t, "mailboxChecked", result["resultDetail"])
	assert.Equal(t, mailck.Valid.Message, result["message"])
}

func Test_MailserverError(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.MailserverError, errors.New("some error")))
	req, err := http.NewRequest("GET", "http://localhost:3000/api/validate?mail=foo@example.com", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, "error", result["result"])
	assert.Equal(t, "mailserverError", result["resultDetail"])
	assert.Equal(t, mailck.MailserverError.Message, result["message"])
}

func Test_ServiceError(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.ServiceError, errors.New("some error")))
	req, err := http.NewRequest("GET", "http://localhost:3000/api/validate?mail=foo@example.com", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, "error", result["result"])
	assert.Equal(t, "serviceError", result["resultDetail"])
	assert.Equal(t, mailck.ServiceError.Message, result["message"])
}

func getJson(t *testing.T, resp *httptest.ResponseRecorder) map[string]interface{} {
	result := map[string]interface{}{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	return result
}
