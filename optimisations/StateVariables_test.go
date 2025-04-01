package optimisations_test

import (
	opt "solsa/optimisations"

	fp "path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateVariableOptimisable(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/not_opt_state.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)
	assert.True(t, opt.StateVariableOptimisable(contracts[0]))

}

func TestStateVariableNonOptimisable(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/opt_state.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)

	assert.False(t, opt.StateVariableOptimisable(contracts[0]))

}

func TestOptimiseStateVariables(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/not_opt_state.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)
	assert.True(t, opt.OptimiseStateVariables(contracts[0]))

}
