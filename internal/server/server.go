package server

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
)

type Server struct {
	Hostname string
	Port     string
	//listener net.Listener
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
		// buffer := make([]byte, 1024)
		// n, _ := conn.Read(buffer)
		// fmt.Println(n)
		// time.Sleep(time.Second * 2)
		// conn.Close()

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	fmt.Printf("Received request from %s", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	buffer := bytes.Buffer{}
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if s == "\r\n" {
			break
		}
		buffer.WriteString(s)
		fmt.Print(s)
	}

	body := "Hello, this is from a valid backend server!"

	// 2. Construct the full HTTP response string.
	// Note the essential `\r\n` line endings.
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n"+ // This empty line is the required separator between headers and body.
			"%s",
		len(body), body,
	)

	// 3. Write the response and close the connection.
	conn.Write([]byte(response))
	conn.Close()
}
