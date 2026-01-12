.PHONY: run build test benchmark
default: run
run:
	go run cmd/s2t/main.go
build:
	go build -o bin/s2t cmd/s2t/main.go
test:
	go test -v -timeout 30s github.com/dragonchen-tw/tongwen-cli-go/cmd/s2t
benchmark:
	go test -v -bench=. -run=none ./cmd/s2t