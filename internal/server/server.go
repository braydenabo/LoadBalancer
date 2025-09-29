package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	Hostname string
	Port     string
}

func (s *Server) Start() {
	fmt.Printf("Starting backend at %s:%s\n", s.Hostname, s.Port)
	ln, err := net.Listen("tcp", s.Hostname+":"+s.Port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Connection from %s\n", conn.RemoteAddr())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close() // Defer connection close, so we always close connection after function exits
	addr := conn.RemoteAddr().String()
	fmt.Printf("Received request from %s", addr)

	reader := bufio.NewReader(conn)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from %s: %v\n", addr, err)
		}
		// End of http headers
		if s == "\r\n" {
			break
		}

	}

	// HTTP response to send back to the load balancer
	body := "Hello, this is from a backend server!"
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n"+
			"%s",
		len(body), body,
	)

	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing to %s: %v\n", addr, err)
	}
}
