package loadbalancer

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type server struct {
	address string
	active  bool
}

type LoadBalancer struct {
	servers        []*server // Backend servers we have access too
	wg             sync.WaitGroup
	current_server int
	mu             sync.Mutex
}

func (lb *LoadBalancer) Start() {
	fmt.Printf("Starting Load Balancer\n")
	lb.current_server = -1
	lb.servers = []*server{
		{address: "127.0.0.1:8081", active: true},
		{address: "127.0.0.1:8082", active: true},
		{address: "127.0.0.1:8083", active: true},
	}

	go lb.healthCheck()

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

	nextServer := lb.getServer()

	backendConn, err := net.Dial("tcp", nextServer.address)
	if err != nil {
		log.Fatal(err)
	}

	go io.Copy(backendConn, conn)
	writer := io.MultiWriter(conn, os.Stdout)

	io.Copy(writer, backendConn)

}

func (lb *LoadBalancer) getServer() *server {
	lb.mu.Lock()
	defer lb.mu.Unlock()
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

func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(time.Second * 10)

	for range ticker.C {
		fmt.Printf("Health check loop\n")

		// Create a map of server health: serverAddress -> health
		// This is to avoid locking our mutex everytime we want to dial and check the server
		healthStatus := make(map[string]bool)
		for _, s := range lb.servers {
			// Dial the server to check its health status
			conn, err := net.DialTimeout("tcp", s.address, 2*time.Second)
			if err != nil {
				healthStatus[s.address] = false
			} else {
				conn.Close()
				healthStatus[s.address] = true

			}
		}

		// Using our server map, we can now use our lock to update server heatlh status
		lb.mu.Lock()
		for _, s := range lb.servers {
			isHealthy := healthStatus[s.address]
			if s.active && !isHealthy {
				log.Printf("Server %s, is down\n", s.address)
				s.active = false
			} else if !s.active && isHealthy {
				log.Printf("Server %s is now running.\n", s.address)
				s.active = true
			}
		}
		// Unlock our mutex
		lb.mu.Unlock()
	}
}
