package mailck

import (
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
)

var emailRexp = regexp.MustCompile("^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}$")

// Check checks the syntax and if valid, it checks the mailbox by connecting to
// the target mailserver
// The fromEmail is used as from address in the communication to the foreign mailserver.
func Check(fromEmail, checkEmail string) (result CheckResult, textMessage string, err error) {
	if !CheckSyntax(checkEmail) {
		return Invalid, "invalid syntax", nil
	}

	if CheckDisposable(checkEmail) {
		return Disposable, "disposable email", nil
	}
	return CheckMailbox(fromEmail, checkEmail)
}

// CheckSyntax returns true for a valid email, false otherwise
func CheckSyntax(checkEmail string) bool {
	return emailRexp.Match([]byte(checkEmail))
}

// CheckMailbox checks the checkEmail by connecting to the target mailbox and returns the result.
// The fromEmail is used as from address in the communication to the foreign mailserver.
func CheckMailbox(fromEmail, checkEmail string) (result CheckResult, textMessage string, err error) {
	mxList, err := net.LookupMX(hostname(checkEmail))
	if err != nil || len(mxList) == 0 {
		return Invalid, "error, no mailserver for hostname", nil
	}
	return checkMailbox(fromEmail, checkEmail, mxList, 25)
}

func checkMailbox(fromEmail, checkEmail string, mxList []*net.MX, port int) (result CheckResult, textMessage string, err error) {
	// try to connect to one mx
	var c *smtp.Client
	for _, mx := range mxList {
		c, err = smtp.Dial(fmt.Sprintf("%v:%v", mx.Host, port))
		if err == nil {
			break
		}
	}
	if err != nil {
		return Undefined, "smtp error", err
	}
	defer c.Close()
	defer c.Quit() // defer ist LIFO

	// HELO
	err = c.Hello(hostname(fromEmail))
	if err != nil {
		return Undefined, "smtp error", err
	}

	// MAIL FROM
	err = c.Mail(fromEmail)
	if err != nil {
		return Undefined, "smtp error", err
	}

	// RCPT TO
	id, err := c.Text.Cmd("RCPT TO:<%s>", checkEmail)
	if err != nil {
		return Undefined, "smtp error", err
	}
	c.Text.StartResponse(id)
	code, msg, err := c.Text.ReadResponse(25)
	c.Text.EndResponse(id)
	if code == 550 {
		return Invalid, "mailbox unavailable", nil
	}

	if err != nil {
		return Undefined, msg, err
	}

	return Valid, "Ok", nil
}

func hostname(mail string) string {
	return mail[strings.Index(mail, "@")+1:]
}
