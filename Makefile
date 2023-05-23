run-local:
		docker-compose up -d
		GO111MODULE=on go mod vendor
		GO111MODULE=on go mod download 
		go run main.go