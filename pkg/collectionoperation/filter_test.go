package collectionoperation

import "testing"

func TestFilter(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := []int{2, 4, 6, 8, 10}

	result := Filter(input, func(i int) bool {
		return i%2 == 0
	})
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], v)
		}
	}
}
