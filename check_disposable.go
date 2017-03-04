package mailck

import (
	"strings"
)

// CheckDisposable returns true if the mail is a disposal mail, false otherwise
func CheckDisposable(checkEmail string) bool {
	host := strings.ToLower(hostname(checkEmail))
	for _, d := range DisposableDomains {
		if host == d {
			return true
		}
	}
	return false
}
