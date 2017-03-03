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
			assert.Equal(t, !test.valid, result.Unverified)
			if !test.valid {
				assert.Equal(t, "invalid mail syntax", result.Msg)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		mail   string
		result CheckResult
	}{
		{"xxx", CheckResult{Unverified: true, Msg: "invalid mail syntax"}},
		{"s.mancke@sdcsdcsdcsdctarent.de", CheckResult{Unverified: true, Msg: "error, no mailserver for hostname"}},
		{"s.mancke@tarent.de", CheckResult{Verified: true, Msg: "Ok"}},
		{"s.mancke+fo42@tarent.de", CheckResult{Verified: true, Msg: "Ok"}},
		{"not_existant@tarent.de", CheckResult{Unverified: true, Msg: "mailbox unavailable"}},
		//{"sebastian@mancke.net", CheckResult{Verified: true, Msg: "Ok"}},
		{"foo@mailinator.com", CheckResult{Verified: true, Msg: "Ok"}},
	}

	for _, test := range tests {
		t.Run(test.mail, func(t *testing.T) {
			start := time.Now()
			result := Check("noreply@mancke.net", test.mail)
			assert.Equal(t, test.result, result)
			fmt.Printf("check for %v: %v\n", test.mail, time.Since(start))
		})
	}
}
