package utils

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

func NormalizeWords(text string) ([]string, error) {
	lowerText := strings.ToLower(text)
	words := strings.FieldsFunc(lowerText, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	output, err := stemWords(words)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func stemWords(words []string) ([]string, error) {
	output := make([]string, 0, len(words))

	for _, word := range words {
		stemmed, err := snowball.Stem(word, "english", false)
		if err != nil {
			return nil, err
		}
		if stemmed == "" {
			continue
		}
		output = append(output, stemmed)
	}

	return output, nil
}
