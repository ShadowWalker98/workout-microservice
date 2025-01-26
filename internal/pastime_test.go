package internal

import "testing"

func TestAddNums(t *testing.T) {
	var tests = []struct {
		name   string
		input1 int
		input2 int
		want   int
	}{
		// the table itself
		{"1 + 1 should be 2", 1, 1, 2},
		{"1 + 2 should be 3", 1, 2, 3},
		{"1 + 3 should be 4", 1, 3, 4},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := addNums(tt.input1, tt.input2)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
