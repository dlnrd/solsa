package optimisations_test

import (
	opt "solsa/optimisations"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariablePacking(t *testing.T) {
	variables := []opt.Variable{{Index: 0, Size: 60}, {Index: 1, Size: 200}, {Index: 2, Size: 33}, {Index: 3, Size: 123}, {Index: 4, Size: 256}, {Index: 5, Size: 82}, {Index: 6, Size: 2}, {Index: 7, Size: 50}, {Index: 8, Size: 159}, {Index: 9, Size: 232}}

	expected := [][]opt.Variable{{{Index: 4, Size: 256}}, {{Index: 9, Size: 232}, {Index: 6, Size: 2}}, {{Index: 1, Size: 200}, {Index: 7, Size: 50}}, {{Index: 8, Size: 159}, {Index: 5, Size: 82}}, {{Index: 3, Size: 123}, {Index: 0, Size: 60}, {Index: 2, Size: 33}}}

	result := opt.VariablePacking(variables)
	assert.Equal(t, expected, result)
}
