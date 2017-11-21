package integration

import (
	"bufio"
	"net"
	"testing"
)

// Counter instance for benchmarks.
var bc = Counter{parallelCount: map[string]int{}, failCount: map[string]int{}}

func BenchmarkLoginCreateAccount(b *testing.B) {
	action := "LoginCreateAccount"
	bc.Parallel(action)
	for i := 0; i < b.N; i++ {
		if conn, err := net.Dial("tcp", *loginAddress); err != nil {
			bc.Fail(action)
		} else {
			if tcpconn, ok := conn.(*net.TCPConn); ok {
				tcpconn.SetLinger(-1)
				tcpconn.SetNoDelay(true)
			}
			buffer := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
			bname := username + bc.Next()
			c := Client{}
			var cm PacketCommand
			if err = c.CreateAccount(buffer, bname, password, CreateAccountMale); err != nil {
				bc.Fail(action)
			} else if err = BinaryRead(buffer, &cm); err != nil || cm != LoginToClientSuccess0069 {
				bc.Fail(action)
			} else if err = c.LoginAccessSuccess(buffer); err != nil {
				bc.Fail(action)
			} else {
				bc.SaveAccount(bname)
			}
			conn.Close()
		}
	}
}

func BenchmarkLoginRequestAccess(b *testing.B) {
	action := "LoginRequestAccess"
	bc.Parallel(action)
	for i := 0; i < b.N; i++ {
		if conn, err := net.Dial("tcp", *loginAddress); err != nil {
			bc.Fail(action)
		} else {
			if tcpconn, ok := conn.(*net.TCPConn); ok {
				tcpconn.SetLinger(-1)
				tcpconn.SetNoDelay(true)
			}
			buffer := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
			c := Client{}
			var cm PacketCommand
			if err = c.RequestAccess(buffer, bc.GetAccount(), password); err != nil {
				bc.Fail(action)
			} else if err = BinaryRead(buffer, &cm); err != nil || cm != LoginToClientSuccess0069 {
				bc.Fail(action)
			} else if err = c.LoginAccessSuccess(buffer); err != nil {
				bc.Fail(action)
			}
			conn.Close()
		}
	}
}

func BenchmarkLoginSummary(b *testing.B) {
	bc.RLock()
	defer bc.RUnlock()
	b.Logf("Benchmark Summary with username (%s) running against address (%s)\n", username, *loginAddress)
	b.Logf("\tCreated Accounts: %d", bc.account)
	b.Logf("\tFailed Creation Requests: %d", bc.failCount["LoginCreateAccount"])
	b.Logf("\tFailed Login Requests: %d", bc.failCount["LoginRequestAccess"])
}
