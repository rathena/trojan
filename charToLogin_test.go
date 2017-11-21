package integration

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

// Create a connection and optimize it for immediate transmission.
//
// Send a Char Server registration command.
//
// Verify valid response packet command and code.
//
// The account id is used to identify each server instance, so a new
// login will be needed for each instance when run against rathena.
func CharToLoginRequestAccess(t *testing.T) {
	var err error
	if charConn, err = net.Dial("tcp", *loginAddress); err != nil {
		t.Fatalf("failed to connect to login server: %s", err)
	}

	// optimize connection
	if tcpconn, ok := charConn.(*net.TCPConn); ok {
		tcpconn.SetLinger(-1)
		tcpconn.SetNoDelay(true)
	}

	// set deadline to avoid blocking tests
	charConn.SetDeadline(time.Now().Add(time.Second * time.Duration(3)))

	// create a buffered readwriter
	buffer := bufio.NewReadWriter(bufio.NewReader(charConn), bufio.NewWriter(charConn))

	// parse charCredentials
	credentials := strings.Split(*charCredentials, ":")
	if len(credentials) < 2 {
		t.Fatalf("charCredentials (%s) are invalid", *charCredentials)
	}

	// prepare properties to send to the login server for a connection request
	var world_username [24]byte
	copy(world_username[:], []byte(credentials[0])[:23])
	var world_password [24]byte
	copy(world_password[:], []byte(credentials[1])[:23])
	var server_name [20]byte = [20]byte{'t', 'e', 's', 't'}
	var server_address uint32 = 0x100007f
	var server_port uint16 = 6122
	var server_type uint16
	var server_new uint16

	// prepare response properties
	var command PacketCommand
	var code uint8

	// Worth mentioning that reading the original source gave me the wrong
	// impression regarding placement of unknown byte chunks, and the offset
	// placement for the server name was equally confusing.
	if err = BinaryWrite(buffer, CharToLoginRegister, world_username, world_password, [4]byte{}, server_address, server_port, server_name, [2]byte{}, server_type, server_new); err != nil {
		t.Fatalf("failed to write char server request: %s", err)
	} else if err = buffer.Flush(); err != nil {
		t.Fatalf("failed to flush char server connection request: %s", err)
	} else if err = BinaryRead(buffer, &command, &code); err != nil {
		t.Fatalf("failed to read login response to char server request: %s", err)
	} else if command != LoginToCharRegistrationReply {
		t.Fatalf("expected command %#04x but got %#04x", LoginToCharRegistrationReply, command)
	} else if code == 3 {
		t.Fatalf("expected code of 0 but got %d", code)
	}
}
