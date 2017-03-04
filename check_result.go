package mailck

// CheckResult contains the information about an email check.
type CheckResult string

var (
	// The mailbox could be completely verified.
	// We know for sure, that the target mailserver would receive mail.
	Valid = CheckResult("valid")

	// Unvalid means, that we know for sure, that the mailbox does not exist.
	// e.g. the server syntax is invalid, the hostname does not exist or the
	// target mailserver said that the mailbox does not exist.
	Invalid = CheckResult("invalid")

	// The mailserver is a throw-away mail gateway like mailinator.com
	Disposable = CheckResult("disposable")

	// Undefined result in case of an error
	Undefined = CheckResult("")
)
