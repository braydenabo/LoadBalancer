# Simple makefile to run loadbalancer, & 3 backend servers

LB_PATH=/cmd/lb
SERVER_PATH=/cmd/server

LB_BINARY=lb
SERVER_BINARY=server

.PHONY: all start build stop clean

start: build
	./$(SERVER_BINARY) -p=8081 &
	./$(SERVER_BINARY) -p=8082 &
	./$(SERVER_BINARY) -p=8083 &

	@sleep 1
	./$(LB_BINARY)

build:
	go build -o $(SERVER_BINARY) .$(SERVER_PATH)
	go build -o $(LB_BINARY) .$(LB_PATH)
	
stop:
	@pkill -f "./$(SERVER_BINARY) -p=" || true
	@pkill -f "./$(LB_BINARY)" || true

clean:
	@rm -f $(LB_BINARY) $(SERVER_BINARY)