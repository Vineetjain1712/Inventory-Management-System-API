# Variables
APP_NAME := inventory-api
PORT := 8080

# Run Go app locally
run:
	go run ./cmd/main.go   

mod:
	go mod tidy 

# Build Go binary locally
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(APP_NAME) ./cmd/main.go


# Run tests
test:
	go test ./...

# Docker build
docker-build:
	docker build -t $(APP_NAME) .

# Docker run
docker-run:
	docker run -p $(PORT):$(PORT) $(APP_NAME)

# Clean binary
clean:
	rm -f $(APP_NAME)
