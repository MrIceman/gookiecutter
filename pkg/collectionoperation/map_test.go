package collectionoperation

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []string{"1", "2", "3", "4", "5"}

	result := Map(input, func(i int) string {
		return fmt.Sprintf("%d", i)
	})

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], v)
		}
	}
}
