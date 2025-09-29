package main

import (
	"loadbalancer/internal/loadbalancer" // Using module name (Correct)
)

// Using module name (Correct)
func main() {
	lb := loadbalancer.LoadBalancer{}
	lb.Start()
}
