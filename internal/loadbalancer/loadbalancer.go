package loadbalancer

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type server struct {
	address string
	active  bool
}

type LoadBalancer struct {
	servers        []*server // Backend servers we have access too
	wg             sync.WaitGroup
	current_server int
}

func (lb *LoadBalancer) Start() {
	fmt.Printf("Starting Load Balancer")
	lb.current_server = -1
	lb.servers = []*server{
		{address: "127.0.0.1:8081", active: true},
		{address: "127.0.0.1:8082", active: true},
		{address: "127.0.0.1:8083", active: true},
	}

	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	lb.acceptRequests(ln)
}

func (lb *LoadBalancer) acceptRequests(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		fmt.Printf("Requests from %s\n", conn.RemoteAddr())
		if err != nil {
			log.Fatal(err)
		}

		lb.wg.Add(1)
		go lb.handleConnections(conn)
	}
}

func (lb *LoadBalancer) handleConnections(conn net.Conn) {
	// Dial a certain server
	fmt.Printf("Received request from %s\n", conn.RemoteAddr())

	//clientResponse := readConn(conn)
	//fmt.Printf("Client %s\n", clientResponse)

	nextServer := lb.getServer()

	backendConn, err := net.Dial("tcp", nextServer.address)
	if err != nil {
		log.Fatal(err)
	}

	// // Write to backend server, the original connection
	// ba, err = backendConn.Write([]byte(clientResponse))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	go io.Copy(backendConn, conn)
	io.Copy(conn, backendConn)

}

// func readConn(conn net.Conn) string {
// 	reader := bufio.NewReader(conn)

// 	buffer := bytes.Buffer{}

// 	for {
// 		s, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if s == "\r\n" {
// 			break
// 		}
// 		buffer.WriteString(s)

// 		if s == "\r\n" {
// 			break
// 		}
// 		fmt.Print(s)
// 	}

// 	return buffer.String()
// }

// Function to get the next server

func (lb *LoadBalancer) getServer() *server {
	next := lb.current_server + 1
	for i := 0; i < len(lb.servers); i++ {
		id := (int(next) + i) % len(lb.servers)
		if lb.servers[id].active {
			lb.current_server = id
			return lb.servers[id]
		}
	}
	return nil
}
