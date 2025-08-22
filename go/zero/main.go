package main

import (
	"log"
	"net"
	"os"

	"github.com/pmoieni/protohackers/go/zero/adapters"
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
			adapters.UseIOCopyBuffer(conn)

			fatal(conn.Close())
		}()
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
