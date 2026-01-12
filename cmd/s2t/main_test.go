package main

import (
	"os"
	"strings"
	"testing"

	"github.com/dragonchen-tw/tongwen-cli-go/internal/assets"
	"github.com/dragonchen-tw/tongwen-cli-go/pkg/converter"
	"github.com/dragonchen-tw/tongwen-cli-go/pkg/loader"
)

func TestLoadJSON(t *testing.T) {
	dict := loader.LoadDiskJSON(assets.Dicts, "s2t-custom.json")

	if dict["“"] != "「" {
		t.Errorf("Expected '“' to map to '「', got '%s'", dict["“"])
	}
}

func TestBuildIndexedMultiMap(t *testing.T) {
	dict := loader.LoadDiskJSON(assets.Dicts, "s2t-phrase.json")
	indexed := converter.MakeMultiIndex(dict)

	if entry, exist := indexed["数字"]; exist {
		if entry.Max != 7 {
			t.Errorf("Expected max key length for index '数字' to be 7, got %d", entry.Max)
		}
		if len(entry.Indies) != 15 {
			t.Errorf("Expected number of entries for index '数字' to be 15, got %d", len(entry.Indies))
		}
		if entry.Indies["数字化"] != "數位化" {
			t.Errorf("Expected '数字化' to map to '數位化', got '%s'", entry.Indies["数字化"])
		}
		if entry.Indies["数字千年版权法"] != "數位千禧年著作權法案" {
			t.Errorf("Expected '数字化' to map to '數位化', got '%s'", entry.Indies["数字化"])
		}
	} else {
		t.Errorf("Index '数字' not found in indexed map")
	}
}

func TestConvertByChar(t *testing.T) {
	dict := loader.LoadDiskJSON(assets.Dicts, "s2t-char.json")
	s2tconverter := converter.NewConverter(dict, nil, false)

	testString := "汉字转换测试123"
	expectedString := "漢字轉換測試123"
	output := s2tconverter.ConvertChar(testString)
	if output != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, output)
	}
}

func TestConvertByPhraseAndChar(t *testing.T) {
	dictChar := loader.LoadDiskJSON(assets.Dicts, "s2t-char.json")
	dictPhrase := loader.LoadDiskJSON(assets.Dicts, "s2t-phrase.json")
	s2tconverter := converter.NewConverter(dictChar, dictPhrase, false)

	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{"余光中的余光有着无限的情怀", "余光中的餘光有著無限的情懷", "Meaningful sentence: should be identity without dict"},
		{"编程语言让世界更美好", "程式語言讓世界更美好", "Meaningful sentence: no changes expected"},
		{"我爱繁体中文与简体中文", "我愛繁體中文與簡體中文", "Mixed Traditional/Simplified: unchanged without mapping"},
	}

	for _, tt := range tests {
		result := s2tconverter.ConvertPhrase(tt.input)
		// fmt.Printf("From: %s\nTo(Char): %s\n", tt.input, s2tconverter.ConvertChar(tt.input))
		// fmt.Printf("To(Phra): %s\n", result)
		if result != tt.expected {
			t.Errorf("%s: expected %s, got %s", tt.name, tt.expected, result)
		}
	}
}

func TestDiffConversionBetweenCharAndPhrase(t *testing.T) {
	// make conversions for char and phrase+char
	// run below test cases
	// 后来她成为了皇后区的皇后 / 後來她成為了皇后區的皇后
	// 树东里巷子里面的面店 / 樹東里巷子裡面的麵店
	// 发芽 头发 江苏 复苏 / 發芽 頭髮 江蘇 復甦
	dictChar := loader.LoadDiskJSON(assets.Dicts, "s2t-char.json")
	dictPhrase := loader.LoadDiskJSON(assets.Dicts, "s2t-phrase.json")
	s2tconverter := converter.NewConverter(dictChar, dictPhrase, false)

	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{"后来她成为了皇后区的皇后", "後來她成為了皇后區的皇后", "Test case 1"},
		{"树东里巷子里面的面店", "樹東里巷子裡面的麵店", "Test case 2"},
		{"发芽 头发 江苏 复苏", "發芽 頭髮 江蘇 復甦", "Test case 33"},
	}

	for _, tt := range tests {
		resultChar := s2tconverter.ConvertChar(tt.input)
		resultPhrase := s2tconverter.ConvertPhrase(tt.input)
		// fmt.Printf("Input: %s\nChar: %s\nPhrase: %s\n", tt.input, resultChar, resultPhrase)
		if resultChar == resultPhrase {
			t.Errorf("%s: expected different results between char and phrase conversion, got same result '%s'", tt.name, resultChar)
		}
		if resultPhrase != tt.expected {
			t.Errorf("%s: expected phrase conversion to be '%s', got '%s'", tt.name, tt.expected, resultPhrase)
		}
	}
}

func BenchmarkTranslateWithCharAndPhrase(b *testing.B) {
	benchmarkData := getBenchmarkData()

	// Here you would typically call your benchmarking functions
	dictChar := loader.LoadDiskJSON(assets.Dicts, "s2t-char.json")
	dictPhrase := loader.LoadDiskJSON(assets.Dicts, "s2t-phrase.json")
	s2tconverter := converter.NewConverter(dictChar, dictPhrase, false)

	b.ResetTimer()
	for range 20 {
		s2tconverter.ConvertPhrase(benchmarkData)
	}
	b.StopTimer()
}

func getBenchmarkData() string {
	data, err := os.ReadFile("benchmark-text.txt")
	if err != nil {
		panic("Failed to read benchmark text file")
	}
	multiplyString := strings.Repeat(string(data), 1000)
	return multiplyString
}
