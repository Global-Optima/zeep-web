package censor

import (
	"encoding/json"
	"io"
	"os"
	"regexp"
	"unicode/utf8"
)

type Regex struct {
	*regexp.Regexp
}

func (r *Regex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	re, err := regexp.Compile(s)
	if err != nil {
		return err
	}
	r.Regexp = re
	return nil
}

type LanguageDictionaries struct {
	Dictionary            []string          `json:"dictionary"`
	FalsePositives        []string          `json:"falsePositives"`
	FalseNegatives        []string          `json:"falseNegatives"`
	CharacterReplacements map[string]string `json:"characterReplacements"`
	ProfanityRegex        Regex             `json:"profanityRegex"`
}

type Dictionaries struct {
	English LanguageDictionaries `json:"en"`
	Russian LanguageDictionaries `json:"ru"`
	Kazakh  LanguageDictionaries `json:"kk"`
}

var dictionaries *Dictionaries

func LoadDictionaries(filename string) (*Dictionaries, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var d Dictionaries
	if err := json.Unmarshal(bytes, &d); err != nil {
		return nil, err
	}

	dictionaries = &d

	return &d, nil
}

func ConvertReplacements(rep map[string]string) map[rune]rune {
	result := make(map[rune]rune)
	for k, v := range rep {
		if len(k) > 0 && len(v) > 0 {
			rKey, _ := utf8.DecodeRuneInString(k)
			rValue, _ := utf8.DecodeRuneInString(v)
			result[rKey] = rValue
		}
	}
	return result
}

func GetDictionaries() *Dictionaries {
	return dictionaries
}
