package server

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
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
	var route string
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from %s: %v\n", addr, err)
		}
		// End of http headers
		if s == "\r\n" {
			break
		}

		// Process the method and request target
		// METHOD_URI_HTTP-Version
		request := strings.Split(s, " ")
		if request[0] == "GET" {
			route = request[1]
		}
	}
	// handle route here
	response := router(route)

	var buf bytes.Buffer
	buf.WriteString(response)
	_, err := conn.Write([]byte(buf.Bytes()))
	if err != nil {
		log.Printf("Error writing to %s: %v\n", addr, err)
	}
}

// router decides what to send, but doesn't worry about formatting.
func router(path string) string {
	switch path {
	case "/":
		body := "Hello, this is from a backend server!"
		return formatHttpResponse(200, "OK", body)

	case "/health":
		// A 204 response has an empty body.
		return formatHttpResponse(204, "No Content", "")

	default:
		// 404 is more accurate for a missing route.
		return formatHttpResponse(404, "Not Found", "404 Not Found")
	}
}

func formatHttpResponse(statusCode int, statusText string, body string) string {
	var responseBuilder strings.Builder

	fmt.Fprintf(&responseBuilder, "HTTP/1.1 %d %s\r\n", statusCode, statusText)
	responseBuilder.WriteString("Connection: close\r\n")

	// Add Content headers if there is a body
	if body != "" {
		fmt.Fprintf(&responseBuilder, "Content-Type: text/plain\r\n")
		fmt.Fprintf(&responseBuilder, "Content-Length: %d\r\n", len(body))
	}
	// End of headers section
	responseBuilder.WriteString("\r\n")

	// Add the body, if it exists
	if body != "" {
		responseBuilder.WriteString(body)
	}
	return responseBuilder.String()
}
