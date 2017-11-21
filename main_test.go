package integration

import (
	"net"
	"testing"
)

// network connections used by various tests
var charConn net.Conn

// Validate all Client to Login Server commands with a clear and concise
// flow control for dependency management.
func TestClientToLogin(t *testing.T) {
	t.Logf("Sending mock Login requests with username %s against login server at %s\n", username, *loginAddress)
	if !t.Run("CreateAccount", ClientToLoginCreateAccount) {
	} else if !t.Run("RequestAccess", ClientToLoginRequestAccess) {
	} else if !t.Run("RequestAccessCommandLatency", ClientToLoginRequestAccessCommandLatency) {
	} else if !t.Run("RequestAccessMessageLatency", ClientToLoginRequestAccessMessageLatency) {
	}
}

// Validate all Char to Login server commands with a clear and
// concise flow control for dependency management.
func TestCharToLogin(t *testing.T) {
	t.Logf("Sending mock world requests with credentials (%s) to login server at %s", *charCredentials, *loginAddress)
	if !t.Run("RequestAccess", CharToLoginRequestAccess) {
		return
	}
}
