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

func testValidationFunction(result mailck.CheckResult, textMessage string, err error) MailValidationFunction {
	return func(checkEmail string) (mailck.CheckResult, string, error) {
		if checkEmail != "foo@example.com" {
			panic("wrong email: " + checkEmail)
		}
		return result, textMessage, err
	}
}

func Test_BadRequest(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.Valid, "Ok", nil))

	req, err := http.NewRequest("POST", "/", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, false, result["success"])
	assert.Equal(t, "missing parameter: email", result["msg"])
}

func Test_SuccessPost(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.Valid, "Ok", nil))
	req, err := http.NewRequest("POST", "http://localhost:3000/api/validate", strings.NewReader(`email=foo@example.com`))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, true, result["success"])
	assert.Equal(t, "valid", result["result"])
	assert.Equal(t, "Ok", result["msg"])
}

func Test_SuccessGet(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.Valid, "Ok", nil))
	req, err := http.NewRequest("GET", "http://localhost:3000/api/validate?email=foo@example.com", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, true, result["success"])
	assert.Equal(t, "valid", result["result"])
	assert.Equal(t, "Ok", result["msg"])
}

func Test_ErrorGet(t *testing.T) {
	handler := NewValidationHandler(testValidationFunction(mailck.Undefined, "smtp error", errors.New("some error")))
	req, err := http.NewRequest("GET", "http://localhost:3000/api/validate?email=foo@example.com", nil)
	assert.NoError(t, err)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

	result := getJson(t, resp)
	assert.Equal(t, false, result["success"])
	assert.Equal(t, "undefined", result["result"])
	assert.Equal(t, "smtp error", result["msg"])
}

func getJson(t *testing.T, resp *httptest.ResponseRecorder) map[string]interface{} {
	result := map[string]interface{}{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	return result
}
