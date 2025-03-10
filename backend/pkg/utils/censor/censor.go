package censor

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

const (
	ALLOWED_SYMBOLS   = `^[a-zA-Zа-яА-ЯёЁәӘіІңҢғҒүҮұҰқҚөӨһҺ0-9\-_ ]+$`
	MAX_WORD_VARIANTS = 100
	MAX_LENGTH        = 50
)

func (c *TextCensorValidator) ValidateText(text string) error {
	nText, err := normalizeText(text)
	if err != nil {
		return err
	}
	if len(nText) < 1 {
		return fmt.Errorf("inappropriate name")
	}
	dicts := GetDictionaries()

	re := regexp.MustCompile(ALLOWED_SYMBOLS)
	if !re.MatchString(text) {
		return fmt.Errorf("regexp is not matched")
	}

	variants, err := dicts.Russian.generateTextVariants(nText)
	if err != nil {
		logrus.Info("ru ", err)
		return err
	}

	for _, variant := range variants {
		if c.RuDetector.IsProfane(variant) {
			logrus.Info("russian variant " + variant)
			return fmt.Errorf("ru: input contains inappropriate words")
		}
	}
	variants, err = dicts.English.generateTextVariants(nText)
	if err != nil {
		logrus.Info("en ", err)
		return err
	}

	for _, variant := range variants {
		if c.EnDetector.IsProfane(variant) {
			logrus.Info("english variant " + variant)
			return fmt.Errorf("en: input contains inappropriate words")
		}
	}

	variants, err = dicts.English.generateTextVariants(nText)
	if err != nil {
		logrus.Info("kk ", err)
		return err
	}

	for _, variant := range variants {
		logrus.Info("kazakh variant " + variant)
		if c.KkDetector.IsProfane(variant) {
			return fmt.Errorf("kk: input contains inappropriate words")
		}
	}

	return nil
}

func normalizeText(text string) (string, error) {
	if len(text) > MAX_LENGTH {
		return "", fmt.Errorf("text too long")
	}
	text = strings.ToLower(strings.ReplaceAll(text, " ", ""))
	re := regexp.MustCompile(`[-_]`)
	return re.ReplaceAllString(text, ""), nil
}
