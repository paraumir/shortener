package service

import (
	"crypto/rand"
	"math/big"
)

const (
	shortURLLength = 10
	alphabet       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func GenerateShortURL() (string, error) {
	result := make([]byte, shortURLLength)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
}

func IsValidShortURL(shortURL string) bool {
	if len(shortURL) != shortURLLength {
		return false
	}

	for _, char := range shortURL {
		if !isAllowedCharacter(char) {
			return false
		}
	}

	return true
}

func isAllowedCharacter(char rune) bool {
	for _, allowed := range alphabet {
		if char == allowed {
			return true
		}
	}

	return false
}
