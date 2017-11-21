package integration

import (
	"bufio"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

// Create a connection to the Login Server and optimize it to send immediately.
//
// Issue an request access command with account creation values.
//
// Verify the response.
func ClientToLoginCreateAccount(t *testing.T) {
	conn, err := net.Dial("tcp", *loginAddress)
	if err != nil {
		t.Logf("connection failed: %s\n", err)
		t.FailNow()
	}
	defer conn.Close()

	// optimize connection
	if tcpconn, ok := conn.(*net.TCPConn); ok {
		tcpconn.SetLinger(-1)
		tcpconn.SetNoDelay(true)
	}

	// set deadline to avoid blocking tests
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(3)))

	// create a buffered readwriter
	buffer := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// prepare a client to generate mock requests
	c := Client{}
	var cm PacketCommand
	if err = c.CreateAccount(buffer, username, password, CreateAccountMale); err != nil {
		t.Fatalf("failed to request access: %s", err)
	} else if err = BinaryRead(buffer, &cm); err != nil {
		t.Fatalf("could not read command reply: %s", err)
	} else if cm != LoginToClientSuccess0069 {
		t.Fatalf("expected 0x0069, but received %#04x", cm)
	} else if err = c.LoginAccessSuccess(buffer); err != nil {
		t.Fatalf("request failed: %s", err)
	}
}

// Create a connection to the Login Server and optimize it to send immediately.
//
// Issue a request access command.
//
// Verify the response.
func ClientToLoginRequestAccess(t *testing.T) {
	conn, err := net.Dial("tcp", *loginAddress)
	if err != nil {
		t.Logf("connection failed: %s\n", err)
		t.FailNow()
	}
	defer conn.Close()

	// optimize connection
	if tcpconn, ok := conn.(*net.TCPConn); ok {
		tcpconn.SetLinger(-1)
		tcpconn.SetNoDelay(true)
	}

	// set deadline to avoid blocking tests
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(3)))

	// create a buffered readwriter
	buffer := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// prepare a client to generate mock requests
	c := Client{}
	var cm PacketCommand
	if err = c.RequestAccess(buffer, username, password); err != nil {
		t.Fatalf("failed to request access: %s", err)
	} else if err = BinaryRead(buffer, &cm); err != nil {
		t.Fatalf("could not read command reply: %s", err)
	} else if cm != LoginToClientSuccess0069 {
		t.Fatalf("expected 0x0069, but received %#04x", cm)
	} else if err = c.LoginAccessSuccess(buffer); err != nil {
		t.Fatalf("request access failed: %s", err)
	}
}

// Create a connection to the Login Server and optimize it to send immediately.
//
// Issue the first byte of a request access command then fake latency
// for a second before sending the remainder of the message.
//
// Verify the response.
func ClientToLoginRequestAccessCommandLatency(t *testing.T) {
	conn, err := net.Dial("tcp", *loginAddress)
	if err != nil {
		t.Logf("connection failed: %s\n", err)
		t.FailNow()
	}
	defer conn.Close()

	// optimize connection
	if tcpconn, ok := conn.(*net.TCPConn); ok {
		tcpconn.SetLinger(-1)
		tcpconn.SetNoDelay(true)
	}

	// set deadline to avoid blocking tests
	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(3)))

	// create a buffered readwriter
	buffer := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// prepare byte arrays for username and password
	netusername := [24]byte{}
	copy(netusername[:], []byte(username)[:23])
	netpassword := [24]byte{}
	copy(netpassword[:], []byte(password)[:23])

	// write the first byte of a command then flush
	if err = binary.Write(conn, binary.LittleEndian, uint8(0x64)); err != nil {
		t.Fatalf("failed to write request access command: %s", err)
	}

	// sleep for a fixed duration
	time.Sleep(1000 * time.Millisecond)

	// write the rest of the data using the buffer
	if err = binary.Write(buffer, binary.LittleEndian, uint8(0x00)); err != nil {
		t.Fatalf("failed to write remainder of the command: %s", err)
	} else if err = binary.Write(buffer, binary.LittleEndian, packetVer); err != nil {
		t.Fatalf("failed to write packet version: %s", err)
	} else if err = binary.Write(buffer, binary.LittleEndian, netusername[:]); err != nil {
		t.Fatalf("failed to write username: %s", err)
	} else if err = binary.Write(buffer, binary.LittleEndian, netpassword[:]); err != nil {
		t.Fatalf("failed to write password: %s", err)
	} else if err = binary.Write(buffer, binary.LittleEndian, ClientType); err != nil {
		t.Fatalf("failed to write client type: %s", err)
	} else if err = buffer.Flush(); err != nil {
		t.Fatalf("failed to flush remaining packet contents: %s", err)
	}

	// verify the response using the client abstraction
	c := Client{}
	var cm PacketCommand
	if err = BinaryRead(buffer, &cm); err != nil {
		t.Fatalf("could not read command reply: %s", err)
	} else if cm != LoginToClientSuccess0069 {
		t.Fatalf("expected 0x0069, but received %#04x", cm)
	} else if err = c.LoginAccessSuccess(buffer); err != nil {
		t.Fatalf("request access failed: %s", err)
	}
}

// Create a connection to the Login Server and optimize it to send immediately.
//
// Write partial request access operation directly to the connection,
// and sleep before sending the rest of the message via buffer.
//
// Verify that the server waited for the stream and responded correctly.
func ClientToLoginRequestAccessMessageLatency(t *testing.T) {
	conn, err := net.Dial("tcp", *loginAddress)
	if err != nil {
		t.Logf("connection failed: %s\n", err)
		t.FailNow()
	}
	defer conn.Close()

	// optimize connection
	if tcpconn, ok := conn.(*net.TCPConn); ok {
		tcpconn.SetLinger(-1)
		tcpconn.SetNoDelay(true)
	}

	// create a buffered readwriter
	buffer := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// prepare byte arrays for username and password
	netusername := [24]byte{}
	copy(netusername[:], []byte(username)[:23])
	netpassword := [24]byte{}
	copy(netpassword[:], []byte(password)[:23])

	// write partial data directly to the connection
	if err = binary.Write(conn, binary.LittleEndian, ClientToLoginRequestAccess0064); err != nil {
		t.Fatalf("failed to write request access command: %s", err)
	} else if err = binary.Write(conn, binary.LittleEndian, packetVer); err != nil {
		t.Fatalf("failed to write packet version: %s", err)
	}

	// sleep for a fixed duration
	time.Sleep(1000 * time.Millisecond)

	// write the rest of the data using the buffer
	if err = binary.Write(buffer, binary.LittleEndian, netusername[:]); err != nil {
		t.Fatalf("failed to write username: %s", err)
	} else if err = binary.Write(buffer, binary.LittleEndian, netpassword[:]); err != nil {
		t.Fatalf("failed to write password: %s", err)
	} else if err = binary.Write(buffer, binary.LittleEndian, ClientType); err != nil {
		t.Fatalf("failed to write client type: %s", err)
	} else if err = buffer.Flush(); err != nil {
		t.Fatalf("failed to flush remaining packet contents: %s", err)
	}

	// verify the response using the client abstraction
	c := Client{}
	var cm PacketCommand
	if err = BinaryRead(buffer, &cm); err != nil {
		t.Fatalf("could not read command reply: %s", err)
	} else if cm != LoginToClientSuccess0069 {
		t.Fatalf("expected 0x0069, but received %#04x", cm)
	} else if err = c.LoginAccessSuccess(buffer); err != nil {
		t.Fatalf("access request failed: %s", err)
	}
}
