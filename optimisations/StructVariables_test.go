package optimisations_test

import (
	opt "solsa/optimisations"

	fp "path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructVariableOptimisable(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/not_opt_struct.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)

	assert.True(t, opt.StructVariableOptimisable(contracts[0]))

}

func TestStructVariableNonOptimisable(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/opt_struct.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)

	assert.False(t, opt.StructVariableOptimisable(contracts[0]))

}

func TestOptimiseStructVariables(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/not_opt_struct.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)

	assert.True(t, opt.OptimiseStructVariables(contracts[0]))

}
