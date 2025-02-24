package optimisations_test

import (
	opt "solsa/optimisations"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariablePacking(t *testing.T) {
	variables := []opt.Variable{opt.Variable{Index: 0, Size: 60}, opt.Variable{Index: 1, Size: 200}, opt.Variable{Index: 2, Size: 33}, opt.Variable{Index: 3, Size: 123}, opt.Variable{Index: 4, Size: 256}, opt.Variable{Index: 5, Size: 82}, opt.Variable{Index: 6, Size: 2}, opt.Variable{Index: 7, Size: 50}, opt.Variable{Index: 8, Size: 159}, opt.Variable{Index: 9, Size: 232}}

	expected := [][]opt.Variable{[]opt.Variable{opt.Variable{Index: 4, Size: 256}}, []opt.Variable{opt.Variable{Index: 9, Size: 232}, opt.Variable{Index: 6, Size: 2}}, []opt.Variable{opt.Variable{Index: 1, Size: 200}, opt.Variable{Index: 7, Size: 50}}, []opt.Variable{opt.Variable{Index: 8, Size: 159}, opt.Variable{Index: 5, Size: 82}}, []opt.Variable{opt.Variable{Index: 3, Size: 123}, opt.Variable{Index: 0, Size: 60}, opt.Variable{Index: 2, Size: 33}}}

	result := opt.VariablePacking(variables)
	assert.Equal(t, expected, result)
}
