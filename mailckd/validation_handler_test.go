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

func Test_Requests(t *testing.T) {
	tests := []struct {
		title              string
		validationFunction MailValidationFunction
		method             string
		requestType        string
		url                string
		body               string
		responseCode       int
		responseType       string
		result             string
		resultDetail       string
		message            string
	}{
		{
			title:              "valid POST example",
			validationFunction: testValidationFunction(mailck.Valid, nil),
			method:             "POST",
			requestType:        "application/x-www-form-urlencoded",
			body:               "mail=foo@example.com",
			responseCode:       200,
			result:             "valid",
			resultDetail:       "mailboxChecked",
			message:            mailck.Valid.Message,
		},
		{
			title:              "valid JSON POST example",
			validationFunction: testValidationFunction(mailck.Valid, nil),
			method:             "POST",
			requestType:        "application/json",
			body:               `{"mail": "foo@example.com"}`,
			responseCode:       200,
			result:             "valid",
			resultDetail:       "mailboxChecked",
			message:            mailck.Valid.Message,
		},
		{
			title:              "valid GET example",
			url:                "/validate?mail=foo%40example.com",
			validationFunction: testValidationFunction(mailck.Valid, nil),
			method:             "GET",
			responseCode:       200,
			result:             "valid",
			resultDetail:       "mailboxChecked",
			message:            mailck.Valid.Message,
		},

		// error cases
		{
			title:              "missing parameter",
			validationFunction: testValidationFunction(mailck.Valid, nil),
			method:             "POST",
			requestType:        "application/x-www-form-urlencoded",
			url:                "/validate",
			body:               "",
			responseCode:       400,
			result:             "error",
			resultDetail:       "clientError",
			message:            "missing parameter: mail",
		},
		{
			title:              "JSON parsing error",
			validationFunction: testValidationFunction(mailck.Valid, nil),
			method:             "POST",
			requestType:        "application/json",
			body:               `{"mail": "foo@example.com`,
			responseCode:       400,
			result:             "error",
			resultDetail:       "clientError",
			message:            "unexpected end of JSON input",
		},
		{
			title:              "method not allowed",
			validationFunction: testValidationFunction(mailck.Valid, nil),
			method:             "PUT",
			requestType:        "application/json",
			body:               `{"mail": "foo@example.com"}`,
			responseCode:       405,
			result:             "error",
			resultDetail:       "clientError",
			message:            "method not allowed",
		},
		{
			title:              "wrong content type",
			validationFunction: testValidationFunction(mailck.Valid, nil),
			method:             "POST",
			requestType:        "text/plain",
			body:               `{"mail": "foo@example.com"}`,
			responseCode:       415,
			result:             "error",
			resultDetail:       "clientError",
			message:            "Unsupported Media Type",
		},
		{
			title:              "service error",
			validationFunction: testValidationFunction(mailck.ServiceError, errors.New("some error")),
			url:                "/validate?mail=foo%40example.com",
			method:             "GET",
			responseCode:       500,
			result:             "error",
			resultDetail:       "serviceError",
			message:            mailck.ServiceError.Message,
		},
		{
			title:              "mailserver error",
			validationFunction: testValidationFunction(mailck.MailserverError, errors.New("some error")),
			url:                "/validate?mail=foo%40example.com",
			method:             "GET",
			responseCode:       502,
			result:             "error",
			resultDetail:       "mailserverError",
			message:            mailck.MailserverError.Message,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			handler := NewValidationHandler(test.validationFunction)
			url := "/validate"
			if test.url != "" {
				url = test.url
			}
			req, err := http.NewRequest(test.method, url, strings.NewReader(test.body))
			if test.requestType != "" {
				req.Header.Set("Content-Type", test.requestType)
			}
			assert.NoError(t, err)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			assert.Equal(t, test.responseCode, resp.Code)
			assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))

			result := getJson(t, resp)
			assert.Equal(t, test.result, result["result"])
			assert.Equal(t, test.resultDetail, result["resultDetail"])
			assert.Equal(t, test.message, result["message"])
		})
	}
}

func getJson(t *testing.T, resp *httptest.ResponseRecorder) map[string]interface{} {
	result := map[string]interface{}{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	return result
}
