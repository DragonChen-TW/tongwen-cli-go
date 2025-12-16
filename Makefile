.PHONY: run build
default: run
run:
	go run cmd/s2t/main.go
build:
	go build cmd/s2t/main.go -o s2t