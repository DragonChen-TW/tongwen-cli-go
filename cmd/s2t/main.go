package main

import (
	"embed"
)

//go:embed tongwen-dict/*.json
var dicts embed.FS

func main() {
}
