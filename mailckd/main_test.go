package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_BasicEndToEnd(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"mailckd", "-host=localhost", "-port=3000", "-text-logging=false"}
	defer func() { os.Args = originalArgs }()

	go main()

	time.Sleep(time.Second)

	r, err := http.Post("http://localhost:3000/api/validate", "application/x-www-form-urlencoded", strings.NewReader(`email=foo@example.com`))
	assert.NoError(t, err)

	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.Unmarshal(b, &result)
	assert.NoError(t, err)

	assert.Equal(t, "invalid", result["result"])
	assert.Equal(t, "error, no mailserver for hostname", result["msg"])
}
