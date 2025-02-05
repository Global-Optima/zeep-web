package censor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"
)

type Regex struct {
	*regexp.Regexp
}

func (r *Regex) UnmarshalJSON(data []byte) error {
	var pattern string
	if err := json.Unmarshal(data, &pattern); err != nil {
		return fmt.Errorf("failed to unmarshal regex pattern: %w", err)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
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
	filePath, err := getDictionaryFilePath(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to determine dictionary file path: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dictionary file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read dictionary file: %w", err)
	}

	var d Dictionaries
	if err := json.Unmarshal(bytes, &d); err != nil {
		return nil, fmt.Errorf("failed to parse dictionary JSON: %w", err)
	}

	dictionaries = &d
	return &d, nil
}

func getDictionaryFilePath(filename string) (string, error) {
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	absPath := filepath.Join(workingDir, filename)
	if _, err := os.Stat(absPath); err == nil {
		return absPath, nil
	}

	return "", errors.New("dictionary file not found")
}

func ConvertReplacements(rep map[string]string) map[rune]rune {
	result := make(map[rune]rune, len(rep))
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
