package mailck

// Result contains the information about an email check.
type Result struct {
	Result       string `json:"result"`
	ResultDetail string `json:"resultDetail"`
	Message      string `json:"message"`
}

var (
	Valid              = Result{"valid", "mailboxChecked", "The email address is valid."}
	InvalidSyntax      = Result{"invalid", "invalidSyntax", "The email format is invalid."}
	InvalidDomain      = Result{"invalid", "invalidDomain", "The email domain does not exist."}
	MailboxUnavailable = Result{"invalid", "mailboxUnavailable", "The email username does not exist."}
	Disposable         = Result{"invalid", "disposable", "The email is a throw-away address."}
	MailserverError    = Result{"error", "mailserverError", "The target mailserver responded with an error."}
	ServiceError       = Result{"error", "serviceError", "An internal error occured while checking."}
)
