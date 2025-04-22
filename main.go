package main

import (
	"solsa/solsa" // good stuff here for exporting to tests :)
)

func main() {
	filePath, _, silent, ok := solsa.ParseFlags()
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

	solsa.OptimiseContracts(builder, silent)

}
