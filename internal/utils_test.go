package internal

import (
	"testing"
)

func TestGetUniqueRandomIntegerFileNames(t *testing.T) {
	LOOP := 100000
	N := 50
	for i := 0; i < LOOP; i++ {
		randomInts := getUniqueRandomIntegerFileNames(1, 100, N)
		uniqueInts := make(map[int]struct{})
		for _, v := range randomInts {
			uniqueInts[v] = struct{}{}
		}
		if len(uniqueInts) != N {
			t.Errorf("expected %d unique random integers, got %v", N, randomInts)
		}
	}
}
