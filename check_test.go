package mailck

import (
	"fmt"
	"github.com/siebenmann/smtpd"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestCheckSyntax(t *testing.T) {
	tests := []struct {
		mail  string
		valid bool
	}{
		{"", false},
		{"xxx", false},
		{"s.mancketarent.de", false},
		{"s.mancke@tarentde", false},
		{"s.mancke@tarent@sdc.de", false},
		{"s.mancke@tarent.de", true},
		{"s.Mancke+yzz42@tarent.de", true},
	}

	for _, test := range tests {
		t.Run(test.mail, func(t *testing.T) {
			result := CheckSyntax(test.mail)
			assert.Equal(t, test.valid, result)
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		mail   string
		result CheckResult
		msg    string
		err    error
	}{
		{"xxx", Invalid, "invalid syntax", nil},
		{"s.mancke@sdcsdcsdcsdctarent.de", Invalid, "error, no mailserver for hostname", nil},
		{"foo@example.com", Invalid, "error, no mailserver for hostname", nil},
		{"foo@mailinator.com", Disposable, "disposable email", nil},
	}

	for _, test := range tests {
		t.Run(test.mail, func(t *testing.T) {
			start := time.Now()
			result, msg, err := Check("noreply@mancke.net", test.mail)
			assert.Equal(t, test.result, result)
			assert.Equal(t, test.msg, msg)
			assert.Equal(t, test.err, err)
			fmt.Printf("check for %30v: %-15v => %-10v (%v)\n", test.mail, time.Since(start), test.result, msg)
		})
	}
}

func Test_checkMailbox(t *testing.T) {
	tests := []struct {
		mail   string
		stopAt smtpd.Command
		result CheckResult
		msg    string
		err    error
	}{
		{"s.mancke@tarent.de", smtpd.QUIT, Valid, "Ok", nil},
		{"not_existant@tarent.de", smtpd.RCPTTO, Invalid, "mailbox unavailable", nil},
	}

	for _, test := range tests {
		t.Run(test.mail, func(t *testing.T) {
			start := time.Now()
			dummyServer := NewDummySMTPServer("localhost:2525", test.stopAt)
			defer dummyServer.Close()
			result, msg, err := checkMailbox("noreply@mancke.net", test.mail, []*net.MX{{Host: "localhost"}}, 2525)
			assert.Equal(t, test.result, result)
			assert.Equal(t, test.msg, msg)
			assert.Equal(t, test.err, err)
			fmt.Printf("check for %30v: %-15v => %-10v (%v)\n", test.mail, time.Since(start), test.result, msg)
		})
	}
}

type DummySMTPServer struct {
	listener net.Listener
	running  bool
	rejectAt smtpd.Command
}

func NewDummySMTPServer(listen string, rejectAt smtpd.Command) *DummySMTPServer {
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		panic(err)
	}
	smtpserver := &DummySMTPServer{
		listener: ln,
		running:  true,
		rejectAt: rejectAt,
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			smtpserver.handleClient(conn)
		}
	}()
	return smtpserver
}

func (smtpserver *DummySMTPServer) Close() {
	smtpserver.listener.Close()
	smtpserver.running = false
}

func (smtpserver *DummySMTPServer) handleClient(conn net.Conn) {
	cfg := smtpd.Config{
		LocalName: "testserver",
		SftName:   "testserver",
	}
	c := smtpd.NewConn(conn, cfg, nil)
	for smtpserver.running {
		event := c.Next()
		//fmt.Printf("event: %+v\n", event)
		if event.Cmd == smtpserver.rejectAt {
			c.Reject()
		} else {
			c.Accept()
		}
		if event.What == smtpd.DONE {
			return
		}
	}
}
