package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net"
	"os"
)

var (
	errMalformedReq = errors.New("malformed request")
)

type req struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

// implement unmarshaler
func (r *req) UnmarshalJSON(bs []byte) error {
	type Alias req

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(bs, &aux); err != nil {
		return err
	}

	if aux.Method != "isPrime" || aux.Number == nil {
		return errMalformedReq
	}

	return nil
}

type res struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func (r *res) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Server struct {
	l net.Listener
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp4", ":"+port)
	fatal(err)

	s.l = listener
	s.listen()
	return nil
}

func (s *Server) Close() error {
	return s.l.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.l.Accept()
		fatal(err)

		go func() {
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				var request req
				if err := json.Unmarshal(scanner.Bytes(), &request); err != nil {
					_, err := conn.Write([]byte("bingus"))
					fatal(err)
					return
				}

				handleReq(conn, &request)
			}
			fatal(scanner.Err())

			defer conn.Close()
		}()
	}
}

func main() {
	s := &Server{}
	s.Run()
}

func handleReq(conn net.Conn, r *req) {
	isInt := math.Trunc(*r.Number) == *r.Number
	response := &res{Method: "isPrime", Prime: false}
	if !isInt {
		bs, err := response.Marshal()
		fatal(err)
		_, err = conn.Write(bs)
		fatal(err)
		return
	}

	isPrime := (func(n int) bool {
		for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
			if n%i == 0 {
				return false
			}
		}
		return n > 1

	})(int(*r.Number))

	if !isPrime {
		bs, err := response.Marshal()
		fatal(err)
		_, err = conn.Write(bs)
		fatal(err)
		return
	}

	response.Prime = true

	bs, err := response.Marshal()
	fatal(err)
	_, err = conn.Write(bs)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
