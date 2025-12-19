package loader

import (
	"embed"
	"encoding/json"
	"log"
)

func LoadDiskJSON(disk embed.FS, filename string) map[string]string {
	filename = "tongwen-dict/" + filename

	data, err := disk.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err)
	}

	return result
}
