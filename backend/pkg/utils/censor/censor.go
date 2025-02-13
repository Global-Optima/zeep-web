package censor

import (
	"fmt"
	"regexp"
	"strings"
)

const ALLOWED_SYMBOLS = `^[a-zA-ZәӘіІңҢғҒүҮұҰқҚөӨһҺ\-_]+$`

func (c *TextCensorValidator) ValidateText(text string) error {
	nText := normalizeText(text)
	dicts := GetDictionaries()

	re := regexp.MustCompile(ALLOWED_SYMBOLS)
	re.MatchString(text)

	if c.RuDetector.IsProfane(nText) {
		return fmt.Errorf("ru: input contains inappropriate words")
	}
	if c.EnDetector.IsProfane(nText) {
		return fmt.Errorf("en: input contains inappropriate words")
	}
	if c.KkDetector.IsProfane(nText) {
		return fmt.Errorf("kk: input contains inappropriate words")
	}

	if containsProfanity(nText, dicts.Russian.ProfanityRegex.Regexp) {
		return fmt.Errorf("ru: input contains inappropriate words (regex)")
	}
	if containsProfanity(nText, dicts.English.ProfanityRegex.Regexp) {
		return fmt.Errorf("en: input contains inappropriate words (regex)")
	}
	if containsProfanity(nText, dicts.Kazakh.ProfanityRegex.Regexp) {
		return fmt.Errorf("kk: input contains inappropriate words (regex)")
	}

	return nil
}

func containsProfanity(text string, pattern *regexp.Regexp) bool {
	return pattern.MatchString(text)
}

func normalizeText(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	re := regexp.MustCompile(`[-_]`)
	return re.ReplaceAllString(text, "")
}
