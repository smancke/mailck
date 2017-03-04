package mailck

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCheckSyntax(t *testing.T) {
	tests := []struct {
		mail  string
		valid bool
	}{
		{"", false},
		{"xxx", false},
		{"s.mancketarent.de", false},
		{"s.mancke@tarentde", false},
		{"s.mancke@tarent@sdc.de", false},
		{"s.mancke@tarent.de", true},
		{"s.Mancke+yzz42@tarent.de", true},
	}

	for _, test := range tests {
		t.Run(test.mail, func(t *testing.T) {
			result := CheckSyntax(test.mail)
			assert.Equal(t, test.valid, result)
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		mail   string
		result CheckResult
		msg    string
		err    error
	}{
		{"xxx", Invalid, "invalid syntax", nil},
		{"s.mancke@sdcsdcsdcsdctarent.de", Invalid, "error, no mailserver for hostname", nil},
		{"foo@example.com", Invalid, "error, no mailserver for hostname", nil},
		{"s.mancke@tarent.de", Valid, "Ok", nil},
		{"s.mancke+fo42@tarent.de", Valid, "Ok", nil},
		{"not_existant@tarent.de", Invalid, "mailbox unavailable", nil},
		//
		//{"sebastian@mancke.net", CheckResult{Valid: true, Msg: "Ok"}},
		{"foo@mailinator.com", Disposable, "disposable email", nil},
	}

	for _, test := range tests {
		t.Run(test.mail, func(t *testing.T) {
			start := time.Now()
			result, msg, err := Check("noreply@mancke.net", test.mail)
			assert.Equal(t, test.result, result)
			assert.Equal(t, test.msg, msg)
			assert.Equal(t, test.err, err)
			fmt.Printf("check for %30v: %-15v => %-10v (%v)\n", test.mail, time.Since(start), test.result, msg)
		})
	}
}
