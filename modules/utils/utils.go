package utils

import (
	"encoding/base64"

	"golang.org/x/exp/constraints"
)

func ContainAny[T constraints.Ordered](src []T, trg T) bool {
	for _, v := range src {
		if v == trg {
			return true
		}
	}

	return false
}

func DecodeToString(v string) (string, error) {
	val, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", err
	}

	return string(val), nil
}
