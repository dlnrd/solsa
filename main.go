package main

import (
	"solsa/solsa" // all the good stuff is here so it can be exported for tests :)
)

func main() {
	filePath, _, ok := solsa.ParseFlags()
	if !ok {
		return
	}

	sources, ok := solsa.GetSources(filePath)
	if !ok {
		return
	}

	builder, ok := solsa.SetupSolgoBuilder(sources)
	if !ok {
		return
	}

	solsa.OptimiseContracts(builder)

}
