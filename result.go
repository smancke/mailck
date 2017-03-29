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
	Valid              = Result{validState.String(), "mailboxChecked", "The email address is valid."}
	InvalidSyntax      = Result{invalidState.String(), "invalidSyntax", "The email format is invalid."}
	InvalidDomain      = Result{invalidState.String(), "invalidDomain", "The email domain does not exist."}
	MailboxUnavailable = Result{invalidState.String(), "mailboxUnavailable", "The email username does not exist."}
	Disposable         = Result{invalidState.String(), "disposable", "The email is a throw-away address."}
	MailserverError    = Result{errorState.String(), "mailserverError", "The target mailserver responded with an error."}
	ServiceError       = Result{errorState.String(), "serviceError", "An internal error occured while checking."}
	clientError        = Result{errorState.String(), "clientError", "The request was was invalid."}
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
