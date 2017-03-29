package mailck

type resultState string

const (
	validState   = resultState("valid")
	invalidState = resultState("invalid")
	errorState   = resultState("error")
)

func (rs resultState) String() string {
	return string(rs)
}

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
	TimeoutError       = Result{"error", "timeoutError", "The connection with the mailserver timed out."}
	NetworkError       = Result{"error", "networkError", "The connection to the mailserver could not be made."}
	ServiceError       = Result{"error", "serviceError", "An internal error occured while checking."}
	clientError        = Result{"error", "clientError", "The request was was invalid."}
)

func (r Result) IsValid() bool {
	return resultState(r.Result) == validState
}

func (r Result) IsInvalid() bool {
	return resultState(r.Result) == invalidState
}

func (r Result) IsError() bool {
	return resultState(r.Result) == errorState
}
