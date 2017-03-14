package mailck

// Result contains the information about an email check.
type Result struct {
	Result       string `json:"result"`
	ResultDetail string `json:"resultDetail"`
	Message      string `json:"message"`
}

var (
	Valid              = Result{"valid", "mailboxChecked", "The mail address is valid."}
	InvalidSyntax      = Result{"invalid", "invalidSyntax", "The format is invalid."}
	InvalidDomain      = Result{"invalid", "invalidDomain", "The domain does not exist."}
	MailboxUnavailable = Result{"invalid", "mailboxUnavailable", "The username does not exist."}
	Disposable         = Result{"invalid", "disposable", "The mailserver is a throw-away mail gateway."}
	MailserverError    = Result{"error", "mailserverError", "The target mailserver responded with an error."}
	ServiceError       = Result{"error", "serviceError", "An internal error occured while checking."}
)
