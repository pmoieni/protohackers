package adapters

import (
	"io"
	"net"
)

func UseIOCopyBuffer(conn net.Conn) error {
	_, err := io.CopyBuffer(conn, conn, nil)
	return err
}

func UseBytes() {

}
