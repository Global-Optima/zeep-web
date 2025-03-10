package censor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Regex struct {
	*regexp.Regexp
}

type CharacterReplacementMap map[string][]string

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
	Dictionary            []string                `json:"dictionary"`
	FalsePositives        []string                `json:"falsePositives"`
	FalseNegatives        []string                `json:"falseNegatives"`
	CharacterReplacements CharacterReplacementMap `json:"characterReplacements"`
}

func (d *LanguageDictionaries) generateTextVariants(text string) ([]string, error) {
	var variants []string
	variants = append(variants, text)

	for key, values := range d.CharacterReplacements {
		var newVariants []string
		for _, variant := range variants {
			if strings.Contains(variant, key) {
				for _, replacement := range values {
					newVariant := strings.ReplaceAll(variant, key, replacement)
					newVariants = append(newVariants, newVariant)
					logrus.Info(newVariant)
				}
			}
		}
		variants = append(variants, newVariants...)
		if len(variants) > MAX_WORD_VARIANTS {
			return nil, errors.New("too many variants to check")
		}
	}

	return variants, nil
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

func GetDictionaries() *Dictionaries {
	return dictionaries
}

func (c *CharacterReplacementMap) UnmarshalJSON(data []byte) error {
	// Create a temporary map for decoding
	temp := make(map[string]interface{})
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to parse character replacements: %w", err)
	}

	// Convert values to []string correctly
	*c = make(map[string][]string)
	for key, value := range temp {
		switch v := value.(type) {
		case string:
			(*c)[key] = []string{v} // Wrap single string into a slice
		case []interface{}:
			var replacements []string
			for _, item := range v {
				if str, ok := item.(string); ok {
					replacements = append(replacements, str)
				} else {
					return fmt.Errorf("invalid character replacement type for key '%s'", key)
				}
			}
			(*c)[key] = replacements
		default:
			return fmt.Errorf("unexpected type for character replacement at key '%s'", key)
		}
	}

	return nil
}
