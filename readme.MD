# Load Balancer
As I attempt to learn go, figured I would make my own load balancer implementation.

```
└── LoadBalancer/
    ├── cmd/
    │   └── main.go (Startup / Configuration)
    ├── internal/
    │   ├── health/
    │   │   └── health.go // Check on server health
    │   ├── server/
    │   │   └── server.go // Backend server setup
    │   ├── loadbalancer/
    │   │   └── loadbalancer.go // Core load balancer logic
    │   └── scheduler/
    │       ├── scheduler.go (Interface for scheduler strategies)
    │       ├── robin.go
    │       └── leastconn.go
    └── go.mod   
 ```

 