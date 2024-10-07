package collectionoperation

import "testing"

func TestGroupBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expected := map[string][]int{
		"odd":  {1, 3, 5, 7, 9},
		"even": {2, 4, 6, 8},
	}

	actual := GroupBy(input, func(i int) string {
		if i%2 == 0 {
			return "even"
		}
		return "odd"
	})

	if len(actual["odd"]) != len(expected["odd"]) {
		t.Errorf("Expected %v, got %v", expected["odd"], actual["odd"])
	}
	if len(actual["even"]) != len(expected["even"]) {
		t.Errorf("Expected %v, got %v", expected["even"], actual["even"])
	}
	for i, v := range actual["odd"] {
		if v != expected["odd"][i] {
			t.Errorf("Expected %v, got %v", expected["odd"], actual["odd"])
		}
	}

	for i, v := range actual["even"] {
		if v != expected["even"][i] {
			t.Errorf("Expected %v, got %v", expected["even"], actual["even"])
		}
	}
}
