package integration

import (
	"encoding/binary"
	"io"
)

// A wrapper to simplify error handling for multiple binary writes.
func BinaryWrite(w io.Writer, d ...interface{}) error {
	var err error
	for a := range d {
		if err = binary.Write(w, binary.LittleEndian, d[a]); err != nil {
			return err
		}
	}
	return nil
}

// A wrapper to simplify error handling for multiple binary reads.
func BinaryRead(r io.Reader, d ...interface{}) error {
	var err error
	for a := range d {
		if err = binary.Read(r, binary.LittleEndian, d[a]); err != nil {
			return err
		}
	}
	return nil
}
