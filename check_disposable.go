package mailck

import (
	"strings"
)

// CheckDisposable returns true if the mail is a disposal mail, false otherwise
func CheckDisposable(checkEmail string) bool {
	host := strings.ToLower(hostname(checkEmail))
	return DisposableDomains[host]
}
