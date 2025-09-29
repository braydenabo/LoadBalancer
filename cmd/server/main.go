package main

import (
	"flag"
	"loadbalancer/internal/server" // Using module name (Correct)
)

func main() {
	var hostname, port string
	flag.StringVar(&hostname, "h", "127.0.0.1", "hostname")
	flag.StringVar(&port, "p", "8081", "port")
	flag.Parse()

	be := server.Server{Hostname: hostname, Port: port}
	be.Start()

}
