package main

import (
	"bytes"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp4", ":"+port)
	fatal(err)
	defer func() {
		fatal(listener.Close())
	}()

	for {
		conn, err := listener.Accept()
		fatal(err)

		go func() {
			var data bytes.Buffer
			_, err = data.ReadFrom(conn)
			fatal(err)

			_, err = conn.Write(data.Bytes())
			fatal(err)

			fatal(conn.Close())
		}()
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
