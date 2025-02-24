package optimisations_test

import (
	"solsa/solsa"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unpackdev/solgo/ir"
)

func SetupTest(t *testing.T, filepath string) *ir.Builder {
	sources, ok := solsa.GetSources(filepath)
	assert.True(t, ok)
	assert.NotNil(t, sources)

	builder, ok := solsa.SetupSolgoBuilder(sources)
	assert.True(t, ok)
	assert.NotNil(t, builder)

	return builder
}
