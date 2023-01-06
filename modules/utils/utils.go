package utils

import (
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
