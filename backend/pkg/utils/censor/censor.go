package censor

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const ALLOWED_SYMBOLS = `^[a-zA-Zа-яА-ЯёЁәӘіІңҢғҒүҮұҰқҚөӨһҺ0-9\-_ ]+$`

func (c *TextCensorValidator) ValidateText(text string) error {
	nText := normalizeText(text)
	if len(nText) < 1 {
		return fmt.Errorf("inappropriate name")
	}
	dicts := GetDictionaries()

	re := regexp.MustCompile(ALLOWED_SYMBOLS)
	if !re.MatchString(text) {
		return fmt.Errorf("regexp is not matched")
	}

	if c.RuDetector.IsProfane(nText) {
		return fmt.Errorf("ru: input contains inappropriate words")
	}
	if c.EnDetector.IsProfane(nText) {
		return fmt.Errorf("en: input contains inappropriate words")
	}
	if c.KkDetector.IsProfane(nText) {
		return fmt.Errorf("kk: input contains inappropriate words")
	}

	containsBadWords, badWords := containsProfanity(nText, dicts.Russian.ProfanityRegex.Regexp)
	if containsBadWords && areFalsePositives(badWords, dicts.Russian.FalsePositives) {
		return fmt.Errorf("ru: input contains inappropriate words (regex)")
	}
	containsBadWords, badWords = containsProfanity(nText, dicts.English.ProfanityRegex.Regexp)
	if containsBadWords && areFalsePositives(badWords, dicts.English.FalsePositives) {
		return fmt.Errorf("en: input contains inappropriate words (regex)")
	}
	containsBadWords, badWords = containsProfanity(nText, dicts.Kazakh.ProfanityRegex.Regexp)
	if containsBadWords && areFalsePositives(badWords, dicts.Kazakh.FalsePositives) {
		return fmt.Errorf("kk: input contains inappropriate words (regex)")
	}

	return nil
}

func areFalsePositives(potentialBadWords, falsePositives []string) bool {
	counter := len(potentialBadWords)
	for _, word := range potentialBadWords {
		if slices.Contains(falsePositives, word) {
			counter--
		}
	}
	return counter == 0
}

func containsProfanity(text string, pattern *regexp.Regexp) (bool, []string) {
	matches := pattern.FindAllString(text, -1)
	if len(matches) > 0 {
		return true, matches
	}
	return false, nil
}

func normalizeText(text string) string {
	text = strings.ToLower(strings.ReplaceAll(text, " ", ""))
	re := regexp.MustCompile(`[-_]`)
	return re.ReplaceAllString(text, "")
}
