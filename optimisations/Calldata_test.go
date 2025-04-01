package optimisations_test

import (
	opt "solsa/optimisations"

	fp "path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalldataOptimisable(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/not_opt_calldata.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)
	assert.True(t, opt.CalldataOptimisable(contracts[0]))

}

func TestCalldataNonOptimisable(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/opt_calldata.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)

	assert.False(t, opt.CalldataOptimisable(contracts[0]))

}

func TestOptimiseCalldataValid(t *testing.T) {
	filepath, _ := fp.Abs("../tests/testdata/not_opt_calldata.sol")
	builder := SetupTest(t, filepath)

	contracts := builder.GetRoot().GetContracts()
	assert.Len(t, contracts, 1)
	assert.True(t, opt.OptimiseCalldata(contracts[0]))
}
