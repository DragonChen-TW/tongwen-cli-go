package main

import (
	"log"
	"os"

	"github.com/dragonchen-tw/tongwen-cli-go/internal/assets"
	"github.com/dragonchen-tw/tongwen-cli-go/pkg/converter"
	"github.com/dragonchen-tw/tongwen-cli-go/pkg/loader"
)

func main() {
	// read input filename, convert the content and save to the same file
	if len(os.Args) < 2 {
		log.Fatal("Usage: s2t <filename>")
	}

	filename := os.Args[1]

	// Read file (UTF-8)
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	content := string(data)

	// Convert content
	converted := convert(content)

	// Write back to the same file (UTF-8)
	err = os.WriteFile(filename, []byte(converted), 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("File converted successfully:", filename)
}

func convert(content string) string {
	dictChar := loader.LoadDiskJSON(assets.Dicts, "s2t-char.json")
	dictPhrase := loader.LoadDiskJSON(assets.Dicts, "s2t-phrase.json")
	s2tconverter := converter.NewConverter(dictChar, dictPhrase, true)
	log.Println("Converter initialized.")
	return s2tconverter.ConvertPhrase(content)
}
