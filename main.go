package main

import (
	"fmt"
	"solsa/solsa" // all the good stuff is here so it can be exported for tests :)
	"time"
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

	startTime := time.Now()
	oState, rState, oStruct, rStruct, oCall, rCall, unopt, failed, totalContracts := solsa.OptimiseContracts(builder, silent)
	timeTaken := time.Now().Sub(startTime)
	sAvg := timeTaken.Seconds() / float64(totalContracts)
	msAvg := sAvg * 1000

	fmt.Println("Total Contacts: ", totalContracts)
	fmt.Printf("Total Time: %.2f\n", timeTaken.Seconds())
	fmt.Println("Average seconds per contract: ", sAvg)
	fmt.Println("Average ms per contract: ", msAvg)
	fmt.Println("Optimised States: ", oState)
	fmt.Println("Optimised Structs: ", oStruct)
	fmt.Println("Optimised Calldata: ", oCall)
	fmt.Println("Unoptimisable: ", unopt)
	fmt.Println("Refactored States: ", rState)
	fmt.Println("Refactored Structs: ", rStruct)
	fmt.Println("Refactored Calldata: ", rCall)
	fmt.Println("Failed Contracts: ", failed)

}
