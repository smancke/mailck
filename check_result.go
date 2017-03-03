package mailck

// CheckResult contains the information about an email check.
type CheckResult struct {
	// The mailbox could be completely verified.
	// We know for sure, that the target mailserver would receive mail.
	Verified bool

	// Unverified means, that we know for sure, that the mailbox does not exist.
	// e.g. the server syntax is invalid, the hostname does not exist or the
	// target mailserver said that the mailbox does not exist.
	Unverified bool

	// The mailserver is a throw-away mail gateway like mailinator.com
	//Disposable bool

	// Err may be a technical error, which occurred while checking.
	// In this case no robust statement about the mail address can be done.
	Err error

	// Detailed Message
	Msg string
}
