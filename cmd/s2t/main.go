package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func LoadLocalJSON(filename string) map[string]string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err)
	}

	return result
}

func main() {
	fmt.Println("s2t converting")

	// result := LoadLocalJSON("tongwen-dict/s2t-char.json")
	result := LoadLocalJSON("tongwen-dict/s2t-phrase.json")

	count := 0
	for k, v := range result {
		fmt.Printf("%s -> %s\n", k, v)

		if count == 5 {
			break
		}
		count++
	}
}
