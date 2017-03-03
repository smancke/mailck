package mailck

import (
	"net"
	"net/smtp"
	"regexp"
	"strings"
)

var emailRexp = regexp.MustCompile("^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}$")

// Check checks the syntax and if valid, it checks the mailbox by connecting to
// the target mailserver
// The fromEmail is used as from address in the communication to the foreign mailserver.
func Check(fromEmail, checkEmail string) (CheckResult, string, error) {
	if !CheckSyntax(checkEmail) {
		return Unvalid, "invalid syntax", nil
	}
	return CheckMailbox(fromEmail, checkEmail)
}

// CheckSyntax returns true for a valid email, false otherwise
func CheckSyntax(checkEmail string) bool {
	return emailRexp.Match([]byte(checkEmail))
}

// CheckMailbox checks the checkEmail by connecting to the target mailbox and returns the result.
// The fromEmail is used as from address in the communication to the foreign mailserver.
func CheckMailbox(fromEmail, checkEmail string) (CheckResult, string, error) {
	mxList, err := net.LookupMX(hostname(checkEmail))
	if err != nil || len(mxList) == 0 {
		return Unvalid, "error, no mailserver for hostname", nil
	}

	var c *smtp.Client
	for _, mx := range mxList {
		c, err = smtp.Dial(mx.Host + ":25")
		if err == nil {
			break
		}
	}
	if err != nil {
		return Undefined, "smtp error", err
	}
	defer c.Close()
	defer c.Quit() // defer ist LIFO

	err = c.Hello(hostname(fromEmail))
	if err != nil {
		return Undefined, "smtp error", err
	}

	err = c.Mail(fromEmail)
	if err != nil {
		return Undefined, "smtp error", err
	}

	id, err := c.Text.Cmd("RCPT TO:<%s>", checkEmail)
	if err != nil {
		return Undefined, "smtp error", err
	}
	c.Text.StartResponse(id)
	code, msg, err := c.Text.ReadResponse(25)
	c.Text.EndResponse(id)
	if code == 550 {
		return Unvalid, "mailbox unavailable", nil
	}

	if err != nil {
		return Undefined, msg, err
	}

	return Valid, "Ok", nil
}

func hostname(mail string) string {
	return mail[strings.Index(mail, "@")+1:]
}
