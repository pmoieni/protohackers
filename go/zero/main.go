package main

import (
	"bytes"
	"fmt"
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
			for {
				data := bytes.NewBuffer(nil)
				buf := make([]byte, 256)
				n, err := conn.Read(buf[0:])
				fatal(err)
				data.Write(buf[0:n])
				fmt.Println(data.String())

				_, err = conn.Write(data.Bytes())
				fatal(err)
			}
		}()
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
