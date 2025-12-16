.PHONY: run build
default: run
run:
	go run cmd/s2t/main.go
build:
	go build -o bin/s2t cmd/s2t/main.go