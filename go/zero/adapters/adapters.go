package adapters

import (
	"bytes"
	"io"
	"net"
)

func UseIOCopyBuffer(conn net.Conn) error {
	_, err := io.CopyBuffer(conn, conn, nil)
	return err
}

func UseBytesBuffer(conn net.Conn) error {
	var bs bytes.Buffer
	if _, err := bs.ReadFrom(conn); err != nil {
		return err
	}

	_, err := conn.Write(bs.Bytes())
	return err
}
