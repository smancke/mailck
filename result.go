package mailck

type resultState string

const (
	validState   resultState = "valid"
	invalidState             = "invalid"
	errorState               = "error"
)

func (rs resultState) String() string {
	return string(rs)
}

// Result contains the information about an email check.
type Result struct {
	Result       resultState `json:"result"`
	ResultDetail string      `json:"resultDetail"`
	Message      string      `json:"message"`
}

var (
	Valid              = Result{validState, "mailboxChecked", "The email address is valid."}
	InvalidSyntax      = Result{invalidState, "invalidSyntax", "The email format is invalid."}
	InvalidDomain      = Result{invalidState, "invalidDomain", "The email domain does not exist."}
	MailboxUnavailable = Result{invalidState, "mailboxUnavailable", "The email username does not exist."}
	Disposable         = Result{invalidState, "disposable", "The email is a throw-away address."}
	MailserverError    = Result{errorState, "mailserverError", "The target mailserver responded with an error."}
	TimeoutError       = Result{errorState, "timeoutError", "The connection with the mailserver timed out."}
	NetworkError       = Result{errorState, "networkError", "The connection to the mailserver could not be made."}
	ServiceError       = Result{errorState, "serviceError", "An internal error occured while checking."}
	clientError        = Result{errorState, "clientError", "The request was was invalid."}
)

func (r Result) IsValid() bool {
	return r.Result == validState
}

func (r Result) IsInvalid() bool {
	return r.Result == invalidState
}

func (r Result) IsError() bool {
	return r.Result == errorState
}
