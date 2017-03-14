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

func Test_ExitOnWrongLoglevel(t *testing.T) {
	exitCode := -1
	osExitOriginal := osExit
	defer func() { osExit = osExitOriginal }()
	osExit = func(code int) {
		exitCode = code
	}
	originalArgs := os.Args
	os.Args = []string{"mailckd", "-log-level=FOOO"}
	defer func() { os.Args = originalArgs }()

	main()
	assert.Equal(t, 1, exitCode)
}

func Test_BasicEndToEnd(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"mailckd", "-host=localhost", "-port=3002", "-text-logging=false"}
	defer func() { os.Args = originalArgs }()

	go main()

	time.Sleep(time.Second)

	r, err := http.Post("http://localhost:3002/api/validate", "application/x-www-form-urlencoded", strings.NewReader(`mail=foo@example.com`))
	assert.NoError(t, err)

	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	result := map[string]interface{}{}
	err = json.Unmarshal(b, &result)
	assert.NoError(t, err)

	assert.Equal(t, "invalid", result["result"])
	assert.Equal(t, "invalidDomain", result["resultDetail"])
	assert.Equal(t, "The email domain does not exist.", result["message"])
}
