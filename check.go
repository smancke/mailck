package mailck

import (
	"context"
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
func Check(fromEmail, checkEmail string) (result Result, err error) {
	if !CheckSyntax(checkEmail) {
		return InvalidSyntax, nil
	}

	if CheckDisposable(checkEmail) {
		return Disposable, nil
	}
	return CheckMailbox(fromEmail, checkEmail)
}

func CheckWithContext(ctx context.Context, fromEmail, checkEmail string) (result Result, err error) {
	if !CheckSyntax(checkEmail) {
		return InvalidSyntax, nil
	}

	if CheckDisposable(checkEmail) {
		return Disposable, nil
	}
	return CheckMailboxWithContext(ctx, fromEmail, checkEmail)
}

// CheckSyntax returns true for a valid email, false otherwise
func CheckSyntax(checkEmail string) bool {
	return emailRexp.Match([]byte(checkEmail))
}

var noContext = context.Background()
var defaultResolver = net.DefaultResolver
var defaultDialer = net.Dialer{}

// CheckMailbox checks the checkEmail by connecting to the target mailbox and returns the result.
// The fromEmail is used as from address in the communication to the foreign mailserver.
func CheckMailbox(fromEmail, checkEmail string) (result Result, err error) {
	mxList, err := net.LookupMX(hostname(checkEmail))
	// TODO: Distinguish between usual network errors
	if err != nil || len(mxList) == 0 {
		return InvalidDomain, nil
	}
	return checkMailbox(noContext, fromEmail, checkEmail, mxList, 25)
}

func CheckMailboxWithContext(ctx context.Context, fromEmail, checkEmail string) (result Result, err error) {
	mxList, err := defaultResolver.LookupMX(ctx, hostname(checkEmail))
	// TODO: Distinguish between usual network errors
	if err != nil || len(mxList) == 0 {
		return InvalidDomain, nil
	}
	return checkMailbox(ctx, fromEmail, checkEmail, mxList, 25)
}

type checkRv struct {
	res Result
	err error
}

func checkMailbox(ctx context.Context, fromEmail, checkEmail string, mxList []*net.MX, port int) (result Result, err error) {
	// try to connect to one mx
	var c *smtp.Client
	for _, mx := range mxList {
		var conn net.Conn
		conn, err = defaultDialer.DialContext(ctx, "tcp", fmt.Sprintf("%v:%v", mx.Host, port))
		if t, ok := err.(*net.OpError); ok {
			if t.Timeout() {
				return TimeoutError, err
			}
			return NetworkError, err
		} else if err != nil {
			return MailserverError, err
		}
		c, err = smtp.NewClient(conn, mx.Host)
		if err == nil {
			break
		}
	}
	if err != nil {
		return MailserverError, err
	}
	if c == nil {
		// just to get very sure, that we have a connection
		// this code line should never be reached!
		return MailserverError, fmt.Errorf("can't obtain connection for %v", checkEmail)
	}

	resChan := make(chan checkRv, 1)

	go func() {
		defer c.Close()
		defer c.Quit() // defer ist LIFO
		// HELO
		// err = c.Hello(hostname(fromEmail))
		err = c.Hello(singleMX(fromEmail))
		if err != nil {
			resChan <- checkRv{MailserverError, err}
			return
		}

		// MAIL FROM
		err = c.Mail(fromEmail)
		if err != nil {
			resChan <- checkRv{MailserverError, err}
			return
		}

		// RCPT TO
		id, err := c.Text.Cmd("RCPT TO:<%s>", checkEmail)
		if err != nil {
			resChan <- checkRv{MailserverError, err}
			return
		}
		c.Text.StartResponse(id)
		code, _, err := c.Text.ReadResponse(25)
		c.Text.EndResponse(id)
		if code == 550 {
			resChan <- checkRv{MailboxUnavailable, nil}
			return
		}

		if err != nil {
			resChan <- checkRv{MailserverError, err}
			return
		}

		resChan <- checkRv{Valid, nil}

	}()
	select {
	case <-ctx.Done():
		return TimeoutError, ctx.Err()
	case q := <-resChan:
		return q.res, q.err
	}
}

func hostname(mail string) string {
	return mail[strings.Index(mail, "@")+1:]
}

func singleMX(email string) string {

	var (
		myList string
		mxLength int
	)

	domain := email[strings.Index(email, "@")+1:]
	mxrecords, _ := net.LookupMX(domain)	
	
	for _, mx := range mxrecords {

		myList = mx.Host
		mxLength = len(myList) - 1
	}
	
	return myList[:mxLength]
}