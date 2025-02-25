package censor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	goaway "github.com/TwiN/go-away"
)

type TextCensorValidator struct {
	EnDetector *goaway.ProfanityDetector
	RuDetector *goaway.ProfanityDetector
	KkDetector *goaway.ProfanityDetector
}

var censorValidator *TextCensorValidator

func initCensorValidator(dicts *Dictionaries) *TextCensorValidator {
	return &TextCensorValidator{
		EnDetector: goaway.NewProfanityDetector().
			WithCustomDictionary(dicts.English.Dictionary, dicts.English.FalsePositives, dicts.English.FalseNegatives).
			WithCustomCharacterReplacements(ConvertReplacements(dicts.English.CharacterReplacements)),
		RuDetector: goaway.NewProfanityDetector().
			WithCustomDictionary(dicts.Russian.Dictionary, dicts.Russian.FalsePositives, dicts.Russian.FalseNegatives).
			WithCustomCharacterReplacements(ConvertReplacements(dicts.Russian.CharacterReplacements)),
		KkDetector: goaway.NewProfanityDetector().
			WithCustomDictionary(dicts.Kazakh.Dictionary, dicts.Kazakh.FalsePositives, dicts.Kazakh.FalseNegatives).
			WithCustomCharacterReplacements(ConvertReplacements(dicts.Kazakh.CharacterReplacements)),
	}
}

func InitializeCensor() error {
	dicts, err := LoadDictionaries("pkg/utils/censor/dictionaries.json")
	if err != nil {
		return err
	}

	censorValidator = initCensorValidator(dicts)
	return nil
}

func GetCensorValidator() *TextCensorValidator {
	if censorValidator == nil {
		panic("censor validator not initialized")
	}
	return censorValidator
}

func InitializeCensorForTests() error {
	baseDir, ok := utils.GetCallerDir(2)
	if !ok {
		return fmt.Errorf("failed to get caller directory")
	}

	candidatePaths := []string{
		"pkg/utils/censor/dictionaries.json",
		"../pkg/utils/censor/dictionaries.json",
		"../../pkg/utils/censor/dictionaries.json",
		"../../../pkg/utils/censor/dictionaries.json",
	}

	var dictPath string
	for _, candidate := range candidatePaths {
		possiblePath := filepath.Join(baseDir, candidate)
		if info, err := os.Stat(possiblePath); err == nil && !info.IsDir() {
			dictPath = filepath.ToSlash(possiblePath)
			break
		}
	}

	if dictPath == "" {
		return fmt.Errorf("could not locate dictionaries.json using candidate paths")
	}

	dicts, err := LoadDictionaries(dictPath)
	if err != nil {
		return fmt.Errorf("failed to load dictionaries from %s: %w", dictPath, err)
	}

	censorValidator = initCensorValidator(dicts)
	return nil
}
