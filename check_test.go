package mailck

import (
	"fmt"
	"github.com/siebenmann/smtpd"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func assertResultState(t *testing.T, result Result, expected resultState) {
	assert.Equal(t, result.IsValid(), expected == validState)
	assert.Equal(t, result.IsInvalid(), expected == invalidState)
	assert.Equal(t, result.IsError(), expected == errorState)
}

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
		mail          string
		result        Result
		err           error
		expectedState resultState
	}{
		{"xxx", InvalidSyntax, nil, invalidState},
		{"s.mancke@sdcsdcsdcsdctarent.de", InvalidDomain, nil, invalidState},
		{"foo@example.com", InvalidDomain, nil, invalidState},
		{"foo@mailinator.com", Disposable, nil, invalidState},
	}

	for _, test := range tests {
		t.Run(test.mail, func(t *testing.T) {
			start := time.Now()
			result, err := Check("noreply@mancke.net", test.mail)
			assert.Equal(t, test.result, result)
			assert.Equal(t, test.err, err)
			assertResultState(t, result, test.expectedState)
			fmt.Printf("check for %30v: %-15v => %-10v (%v)\n", test.mail, time.Since(start), test.result.Result, test.result.ResultDetail)
		})
	}
}

func Test_checkMailbox(t *testing.T) {
	tests := []struct {
		stopAt        smtpd.Command
		result        Result
		expectError   bool
		expectedState resultState
	}{
		{smtpd.QUIT, Valid, false, validState},
		{smtpd.RCPTTO, MailboxUnavailable, false, invalidState},
		{smtpd.MAILFROM, MailserverError, true, errorState},
		{smtpd.HELO, MailserverError, true, errorState},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("stop at: %v", test.stopAt), func(t *testing.T) {
			dummyServer := NewDummySMTPServer("localhost:2525", test.stopAt)
			defer dummyServer.Close()
			result, err := checkMailbox("noreply@mancke.net", "foo@bar.de", []*net.MX{{Host: "localhost"}}, 2525)
			assert.Equal(t, test.result, result)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assertResultState(t, result, test.expectedState)
		})
	}
}

func Test_checkMailbox_NetworkError(t *testing.T) {
	result, err := checkMailbox("noreply@mancke.net", "foo@bar.de", []*net.MX{{Host: "localhost"}}, 6666)
	assert.Equal(t, MailserverError, result)
	assert.Error(t, err)
	assertResultState(t, result, errorState)
}

type DummySMTPServer struct {
	listener net.Listener
	running  bool
	rejectAt smtpd.Command
}

func NewDummySMTPServer(listen string, rejectAt smtpd.Command) *DummySMTPServer {
	time.Sleep(10 * time.Millisecond)
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
	time.Sleep(10 * time.Millisecond)
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
		if event.Cmd == smtpserver.rejectAt ||
			(smtpserver.rejectAt == smtpd.HELO && event.Cmd == smtpd.EHLO) {
			c.Reject()
		} else {
			c.Accept()
		}
		if event.What == smtpd.DONE {
			return
		}
	}
}
