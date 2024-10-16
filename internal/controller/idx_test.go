package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Idx_Can_Be_Negative(t *testing.T) {
	idx, err := FromInt(-1)

	assert.Empty(t, idx)
	assert.ErrorContains(t, err, "less than zero")
}

func Test_Idx_Should_Be_Positive(t *testing.T) {
	idx, err := FromInt(1)

	assert.Equal(t, idx.Value(), 1)
	assert.NoError(t, err)
}

func Test_Idx_Should_Be_In(t *testing.T) {
	var tests = []struct {
		name  string
		index int
		max   int
		want  bool
	}{
		{"0 in [0,1)", 0, 1, true},
		{"1 in [1,1)", 1, 1, false},
		{"2 not it [0,1)", 2, 1, false},
		{"3 in [0,1,2,3)", 3, 3, false},
		{"2 in [0,1,2,3)", 2, 3, true},
		{"4 not in [0,1,2,3,4)", 4, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, _ := FromInt(tt.index)
			result := idx.IsIn(tt.max)
			if result != tt.want {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}
