test:
	go test -v ./...
build:
	CGO_ENABLED=0 go build -o ht ./cmd/http/main.go
