package converter

type DicMap map[string]string

type Chars []rune

type IndexedMap struct {
	Max    int
	Indies DicMap
}

type IndexedMultiMap map[string]*IndexedMap

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

func ConvertChar(dict DicMap, inputed string) string {
	converted := Chars(inputed)

	for i, char := range converted {
		if mapped, exist := dict[string(char)]; exist {
			converted[i] = []rune(mapped)[0]
		}
	}
	return string(converted)
}

func ConvertPhraseAndChar(single DicMap, multi IndexedMultiMap, text string) string {
	converted := ""
	textChars := Chars(text)
	textLength := len(textChars)

	for pointer := 0; pointer < textLength; pointer++ {
		index := textChars[pointer:]
		if len(index) > 1 {
			index = textChars[pointer : pointer+2]
		}
		indexedData, found := multi[string(index)]
		isFound := false

		if found {
			sliceLength := min(textLength-pointer, indexedData.Max)

			for sliceLength > 1 {
				tomap := string(textChars[pointer : pointer+sliceLength])
				if mapped, exists := indexedData.Indies[tomap]; exists {
					converted += mapped
					pointer += sliceLength - 1
					isFound = true
					break
				}
				sliceLength--
			}

			if !isFound {
				mapped, exist := single[string(textChars[pointer])]
				if exist {
					converted += mapped
				} else {
					converted += string(textChars[pointer])
				}
			}
		} else {
			mapped, exist := single[string(textChars[pointer])]
			if exist {
				converted += mapped
			} else {
				converted += string(textChars[pointer])
			}
		}
	}

	return converted
}
