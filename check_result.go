package mailck

// CheckResult contains the information about an email check.
type CheckResult string

// The mailbox could be completely verified.
// We know for sure, that the target mailserver would receive mail.
var Valid = CheckResult("valid")

// Unvalid means, that we know for sure, that the mailbox does not exist.
// e.g. the server syntax is invalid, the hostname does not exist or the
// target mailserver said that the mailbox does not exist.
var Unvalid = CheckResult("unvalid")

// The mailserver is a throw-away mail gateway like mailinator.com
var Disposable = CheckResult("Disposable")

// Undefined result in case of an error
var Undefined = CheckResult("")
