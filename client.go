package integration

import (
	"bufio"
)

// Emulate requests that a client would make.
//
// Collect data from responses for further validation.
type Client struct {
	Account
	Worlds []World
}

// Uses little endian binary encoding to request access with a username and
// password.
func (c *Client) RequestAccess(rw *bufio.ReadWriter, username, password string) error {
	netusername := [24]byte{}
	copy(netusername[:], []byte(username)[:23])
	netpassword := [24]byte{}
	copy(netpassword[:], []byte(password)[:23])
	if err := BinaryWrite(rw, ClientToLoginRequestAccess0064, packetVer, netusername[:], netpassword[:], ClientType); err != nil {
		return err
	}
	return rw.Flush()
}

// Appends the gender suffix to the username and runs RequestAccess.
func (c *Client) CreateAccount(rw *bufio.ReadWriter, username, password string, gender GenderSuffix) error {
	return c.RequestAccess(rw, username+string(gender), password)
}

// Parse a successful login response.
//
// Expects the bytes for the properties and at least one server.
func (c *Client) LoginAccessSuccess(rw *bufio.ReadWriter) error {
	var packetSize uint16
	if _, err := rw.Peek(45 + 32); err != nil {
		return err
	} else if err = BinaryRead(rw, &packetSize, &c.AuthKeyOne, &c.AccountID, &c.AuthKeyTwo); err != nil {
		return err
	} else if _, err = rw.Discard(30); err != nil {
		return err
	} else if err = BinaryRead(rw, &c.Gender); err != nil {
		return err
	}
	c.Worlds = []World{}
	for i := 0; i < (int(packetSize)-45)/32; i++ {
		w := World{}
		if err := BinaryRead(rw, &w.IP, &w.Port, &w.Name, &w.Users, &w.Type, &w.New); err != nil {
			return err
		}
		c.Worlds = append(c.Worlds, w)
	}
	return nil
}
