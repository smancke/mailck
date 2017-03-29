package mailck

import (
	"fmt"
	"github.com/siebenmann/smtpd"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
	"context"
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
		result Result
		err    error
	}{
		{"xxx", InvalidSyntax, nil},
		{"s.mancke@sdcsdcsdcsdctarent.de", InvalidDomain, nil},
		{"foo@example.com", InvalidDomain, nil},
		{"foo@mailinator.com", Disposable, nil},
	}

	for _, test := range tests {
		t.Run(test.mail, func(t *testing.T) {
			start := time.Now()
			result, err := Check("noreply@mancke.net", test.mail)
			assert.Equal(t, test.result, result)
			assert.Equal(t, test.err, err)
			fmt.Printf("check for %30v: %-15v => %-10v (%v)\n", test.mail, time.Since(start), test.result.Result, test.result.ResultDetail)
		})
	}
}

func Test_checkMailbox(t *testing.T) {
	tests := []struct {
		stopAt      smtpd.Command
		result      Result
		expectError bool
	}{
		{smtpd.QUIT, Valid, false},
		{smtpd.RCPTTO, MailboxUnavailable, false},
		{smtpd.MAILFROM, MailserverError, true},
		{smtpd.HELO, MailserverError, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("stop at: %v", test.stopAt), func(t *testing.T) {
			dummyServer := NewDummySMTPServer("localhost:2525", test.stopAt, 0)
			defer dummyServer.Close()
			result, err := checkMailbox(noContext,"noreply@mancke.net", "foo@bar.de", []*net.MX{{Host: "localhost"}}, 2525)
			assert.Equal(t, test.result, result)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_checkMailbox_NetworkError(t *testing.T) {
	result, err := checkMailbox(noContext,"noreply@mancke.net", "foo@bar.de", []*net.MX{{Host: "localhost"}}, 6666)
	assert.Equal(t, NetworkError, result)
	assert.Error(t, err)
}

func Test_checkMailboxContext(t *testing.T) {
	deltas := []struct{
		delayTime      time.Duration
		contextTime    time.Duration
		expectedResult Result
	}{
		{ 0, 0, TimeoutError },
		{ 0, time.Second, Valid },
		{ time.Millisecond * 1500, 200 * time.Millisecond, TimeoutError },
	}
	for _,d := range deltas {
		t.Run(fmt.Sprintf("context time %v delay %v expected %v", d.contextTime,d.delayTime,d.expectedResult.Result), func(t *testing.T) {
			dummyServer := NewDummySMTPServer("localhost:2528", smtpd.QUIT,d.delayTime)
			tt := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), d.contextTime)
			result,err := checkMailbox(ctx, "noreply@mancke.net", "foo@bar.de", []*net.MX{{Host: "127.0.0.1"}}, 2528)
			if d.expectedResult == Valid {
				assert.NoError(t,err)
			} else {
				assert.Error(t,err)
			}
			assert.Equal(t,d.expectedResult, result)
			// confirm that we completed within requested time
			// add 10ms of wiggle room
			assert.WithinDuration(t, time.Now(), tt, d.contextTime + 10 * time.Millisecond)
			dummyServer.Close()
			cancel()
		})
	}
}

type DummySMTPServer struct {
	listener net.Listener
	running  bool
	rejectAt smtpd.Command
	delay time.Duration
}

func NewDummySMTPServer(listen string, rejectAt smtpd.Command, delay time.Duration) *DummySMTPServer {
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Millisecond)
	smtpserver := &DummySMTPServer{
		listener: ln,
		running:  true,
		rejectAt: rejectAt,
		delay: delay,
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go smtpserver.handleClient(conn)
		}
	}()
	time.Sleep(10 * time.Millisecond)
	return smtpserver
}

func (smtpserver *DummySMTPServer) Close() {
	smtpserver.listener.Close()
	smtpserver.running = false
	time.Sleep(10 * time.Millisecond)
}

func (smtpserver *DummySMTPServer) handleClient(conn net.Conn) {
	cfg := smtpd.Config{
		LocalName: "testserver",
		SftName:   "testserver",
	}
	c := smtpd.NewConn(conn, cfg, nil)
	for smtpserver.running {
		event := c.Next()
		time.Sleep(smtpserver.delay)
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
