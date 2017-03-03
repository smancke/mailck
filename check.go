package mailck

import (
	"net"
	"net/smtp"
	"regexp"
	"strings"
)

var emailRexp *regexp.Regexp

func init() {
	emailRexp = regexp.MustCompile("^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}$")
}

// Check checks the syntax and if valid, it checks the mailbox by connecting to
// the target mailserver
// The fromEmail is used as from address in the communication to the foreign mailserver.
func Check(fromEmail, checkEmail string) CheckResult {
	if r := CheckSyntax(checkEmail); r.Unverified {
		return r
	}
	return CheckMailbox(fromEmail, checkEmail)
}

// CheckSyntax verifies, that the syntax of an email is valid.
func CheckSyntax(checkEmail string) CheckResult {
	if !emailRexp.Match([]byte(checkEmail)) {
		return CheckResult{
			Unverified: true,
			Msg:        "invalid mail syntax",
		}
	}
	return CheckResult{}
}

// CheckMailbox checks the checkEmail by connecting to the target mailbox and returns the result.
// The fromEmail is used as from address in the communication to the foreign mailserver.
func CheckMailbox(fromEmail, checkEmail string) CheckResult {
	mxList, err := net.LookupMX(hostname(checkEmail))
	if err != nil || len(mxList) == 0 {
		return CheckResult{
			Unverified: true,
			Msg:        "error, no mailserver for hostname",
		}
	}

	var c *smtp.Client
	for _, mx := range mxList {
		c, err = smtp.Dial(mx.Host + ":25")
		if err == nil {
			break
		}
	}
	if err != nil {
		return CheckResult{Err: err, Msg: "error connecting mailserver"}
	}
	defer c.Close()
	defer c.Quit() // defer ist LIFO

	err = c.Hello(hostname(fromEmail))
	if err != nil {
		return CheckResult{Err: err, Msg: "error on helo with mailserver"}
	}

	err = c.Mail(fromEmail)
	if err != nil {
		return CheckResult{Err: err, Msg: "sender rejected by mailserver"}
	}

	id, err := c.Text.Cmd("RCPT TO:<%s>", checkEmail)
	if err != nil {
		return CheckResult{Err: err, Msg: "communication error on RCPT TO"}
	}
	c.Text.StartResponse(id)
	code, msg, err := c.Text.ReadResponse(25)
	c.Text.EndResponse(id)
	if code == 550 {
		return CheckResult{Unverified: true, Msg: "mailbox unavailable"}
	}

	if err != nil {
		return CheckResult{Err: err, Msg: msg}
	}

	return CheckResult{Verified: true, Msg: "Ok"}
}

func hostname(mail string) string {
	return mail[strings.Index(mail, "@")+1:]
}
