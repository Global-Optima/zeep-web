package censor

import (
	goaway "github.com/TwiN/go-away"
)

type TextCensorValidator struct {
	EnDetector *goaway.ProfanityDetector
	RuDetector *goaway.ProfanityDetector
	KkDetector *goaway.ProfanityDetector
}

var censorValidator *TextCensorValidator

func InitializeCensor() error {

	dicts, err := LoadDictionaries("pkg/utils/censor/dictionaries.json")
	if err != nil {
		return err
	}

	censorValidator = &TextCensorValidator{}

	censorValidator.EnDetector = goaway.NewProfanityDetector().
		WithCustomDictionary(dicts.English.Dictionary, dicts.English.FalsePositives, dicts.English.FalseNegatives).
		WithCustomCharacterReplacements(ConvertReplacements(dicts.English.CharacterReplacements))

	censorValidator.RuDetector = goaway.NewProfanityDetector().
		WithCustomDictionary(dicts.Russian.Dictionary, dicts.Russian.FalsePositives, dicts.Russian.FalseNegatives).
		WithCustomCharacterReplacements(ConvertReplacements(dicts.Russian.CharacterReplacements))

	censorValidator.KkDetector = goaway.NewProfanityDetector().
		WithCustomDictionary(dicts.Kazakh.Dictionary, dicts.Kazakh.FalsePositives, dicts.Kazakh.FalseNegatives).
		WithCustomCharacterReplacements(ConvertReplacements(dicts.Kazakh.CharacterReplacements))

	return nil
}

func GetCensorValidator() *TextCensorValidator {
	if censorValidator == nil {
		panic("censor validator not initialized")
	}
	return censorValidator
}
