package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Benchmarking Tongwen Converter...")
}

func GetBenchmarkData() string {
	data, _ := os.ReadFile("./benchmark-text.txt")
	return string(data)
}
