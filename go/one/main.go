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
	Method *string  `json:"method"`
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

	if aux.Method == nil || *aux.Method != "isPrime" || aux.Number == nil {
		return errMalformedReq
	}

	return nil
}

type res struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func (r *res) Marshal() ([]byte, error) {
	bs, err := json.Marshal(r)
	bs = append(bs, []byte("\n")...)

	return bs, err
}

type Server struct {
	l net.Listener
}

func (s *Server) Run() {
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp4", ":"+port)
	fatal(err)

	s.l = listener
	s.listen()
}

func (s *Server) Close() error {
	return s.l.Close()
}

func (s *Server) listen() error {
	for {
		conn, err := s.l.Accept()
		fatal(err)

		go func() {
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				var request req
				log.Println(scanner.Text())
				if err := json.Unmarshal(scanner.Bytes(), &request); err != nil {
					_, err := conn.Write([]byte("bingus"))
					fatal(err)
					defer conn.Close()
					return
				}

				handleReq(conn, &request)
			}
			fatal(scanner.Err())
		}()
	}
}

func main() {
	s := &Server{}
	s.Run()
	defer s.Close()
}

func handleReq(conn net.Conn, r *req) {
	isInt := math.Trunc(*r.Number) == *r.Number
	response := &res{Method: "isPrime", Prime: false}
	if !isInt {
		log.Println("NOT INT")
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
		log.Println("NOT PRIME")
		bs, err := response.Marshal()
		fatal(err)
		_, err = conn.Write(bs)
		fatal(err)
		return
	}

	log.Println("IT'S PRIME!!!")
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
