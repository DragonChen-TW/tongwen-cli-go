package converter

import (
	"fmt"
	"log"
	"os"
	"time"
)

type DicMap map[string]string

type Chars []rune

type IndexedMap struct {
	Max    int
	Indies DicMap
}

type IndexedMultiMap map[string]*IndexedMap

type S2TConverter struct {
	DictChar   DicMap
	DictPhrase IndexedMultiMap
	Verbose    bool
}

// newConverter creates a new S2TConverter instance that input dictChar and dictPhrase. then do neccesary process
func NewConverter(dictChar DicMap, dictPhrase DicMap, verbose bool) *S2TConverter {
	indexedPhrase := MakeMultiIndex(dictPhrase)
	return &S2TConverter{
		DictChar:   dictChar,
		DictPhrase: indexedPhrase,
		Verbose:    verbose,
	}
}

func MakeMultiIndex(multi DicMap) IndexedMultiMap {
	indexed := make(IndexedMultiMap)
	for key, value := range multi {
		index := ""
		if len(key) >= 2 {
			index = string(Chars(key)[:2])
		} else {
			index = key
		}

		if _, exist := indexed[index]; !exist {
			indexed[index] = &IndexedMap{
				Max:    0,
				Indies: make(map[string]string),
			}
		}

		if len(Chars(key)) > indexed[index].Max {
			indexed[index].Max = len(Chars(key))
		}

		indexed[index].Indies[key] = value
	}
	return indexed
}

func (s *S2TConverter) ConvertChar(inputed string) string {
	converted := Chars(inputed)

	for i, char := range converted {
		if mapped, exist := s.DictChar[string(char)]; exist {
			converted[i] = []rune(mapped)[0]
		}
	}
	return string(converted)
}

func (s *S2TConverter) ConvertPhrase(text string) string {
	startTime := time.Now()
	converted := ""
	textChars := Chars(text)
	textLength := len(textChars)
	if s.Verbose {
		log.Printf("Total length of textChars: %d\n", textLength)
	}

	// Simple progress bar (flush in terminal)
	progressBarWidth := 50
	for pointer := 0; pointer < textLength; pointer++ {
		if s.Verbose {
			if textLength > 0 && pointer%(textLength/progressBarWidth+1) == 0 {
				progress := int(float64(pointer) / float64(textLength) * float64(progressBarWidth))
				bar := "["
				for i := 0; i < progressBarWidth; i++ {
					if i < progress {
						bar += "="
					} else {
						bar += " "
					}
				}
				bar += "]"
				// Print progress bar with carriage return and flush (bash)
				fmt.Fprintf(os.Stdout, "\rProgress: %s %d/%d", bar, pointer, textLength)
				os.Stdout.Sync()
			}
		}

		var index string
		if pointer+2 < textLength {
			index = string(textChars[pointer : pointer+2])
		} else {
			index = string(textChars[pointer:])
		}
		indexedData, found := s.DictPhrase[index]
		isFound := false

		if found {
			sliceLength := min(textLength-pointer, indexedData.Max)

			for ; sliceLength > 1; sliceLength-- {
				tomap := string(textChars[pointer : pointer+sliceLength])
				if mapped, exists := indexedData.Indies[tomap]; exists {
					converted += mapped
					pointer += sliceLength - 1
					isFound = true
					break
				}
			}

			if !isFound {
				mapped, exist := s.DictChar[string(textChars[pointer])]
				if exist {
					converted += mapped
				} else {
					converted += string(textChars[pointer])
				}
			}
		} else {
			mapped, exist := s.DictChar[string(textChars[pointer])]
			if exist {
				converted += mapped
			} else {
				converted += string(textChars[pointer])
			}
		}
	}

	if s.Verbose {
		// Print newline after progress bar is done
		fmt.Fprintln(os.Stdout)

		timeCost := time.Since(startTime)
		log.Printf("ConvertPhrase time cost: %v\n", timeCost)
	}

	return converted
}
