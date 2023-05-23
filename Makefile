run-local:
    docker-compose up -d
	go mod vendor
	go mod download 
	go run main.go