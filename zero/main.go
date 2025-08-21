package main

import (
	"bytes"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp4", ":8080")
	fatal(err)
	defer func() {
		fatal(listener.Close())
	}()

	for {
		conn, err := listener.Accept()
		fatal(err)

		go func() {
			data := bytes.NewBuffer(nil)
			buf := make([]byte, 256)
			n, err := conn.Read(buf[0:])
			fatal(err)
			data.Write(buf[0:n])

			_, err = conn.Write(data.Bytes())
			fatal(err)
		}()
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
