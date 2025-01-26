package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/ir"
)

const maxSlotSize = 256 // in bits

/*
 TODO:
 - find size of the contracts variables
 - use bin packing algorithm to sort into more efficiently packed slots
 - rewrite bin packing algo to potentially find better solutions (less bins)
 - output results
 - write more optimisations :)
*/

func main() {
	tmpFilePath := flag.String("i", "", "Path to File/Directory")
	outputPath := flag.String("o", "", "Path to save output")
	flag.Parse()

	filePath, _ := fp.Abs(*tmpFilePath)

	var isFile bool

	// CHECKS PATH ISNT EMPTY
	if filePath == "" {
		fmt.Println("Enter a path to file or directory (-i)")
		return // better than os.Exit(0) as runs deconstructors that close stuff
	}

	// CHECK INPUT PATH VALIDITY
	if fileInfo, err := os.Stat(filePath); err == nil { // file exists
		isFile = !fileInfo.IsDir()
	} else {
		fmt.Println("File doesn't exist, please enter a valid filepath")
		return
	}
	if *outputPath == "" { // don't want output
		// don't do anything :)
	} else if fileInfo, err := os.Stat(*outputPath); err == nil { // file exists
		if fileInfo.IsDir() {
			fmt.Println("Output file is a directory, please enter a filepath to save")
			return
		}
	} else { // file doesn't exist
		data := []byte("data")
		if err := os.WriteFile(*outputPath, data, 0644); err == nil {
			os.Remove(*outputPath) // can write to outputfile
		} else {
			fmt.Println("Can't write to output file, do you have write permissions?")
			return
		}
	}

	if isFile {
		if fp.Ext(filePath) != ".sol" {
			fmt.Println("File is not a .sol file")
			return
		}
	} else {
		// get a list of all the files to be checked?
	}

	contractName := strings.TrimSuffix(fp.Base(filePath), ".sol")
	sourcesPath := fp.Dir(filePath)

	solgo.SetLocalSourcesPath(sourcesPath)
	sources, err := solgo.NewSourcesFromPath(contractName, sourcesPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	builder, err := ir.NewBuilderFromSources(ctx, sources)
	if err != nil {
		fmt.Println("builder init: ", err)
		return
	}

	if err := builder.Parse(); err != nil {
		fmt.Println("builder parse: ", err)
	}

	if err := builder.Build(); err != nil {
		fmt.Println("builder build: ", err)
	}

	ast := builder.GetAstBuilder()
	if err := ast.ResolveReferences(); err != nil {
		fmt.Println("AST Resolve References: ", err)
	}

	/*
		 TODO: state vaiable sorting
			- identify?
			- calculate slots used
			- calculate potential slots used
				- sum of sv sizes / slot size
	*/

	contracts := builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		fmt.Println("Contract: ", contract.GetName())
		fmt.Println("StateVariableOptimisable: ", StateVariableOptimisable(contract))
	}
}

func StateVariableOptimisable(contract *ir.Contract) bool {
	stateVariables := contract.GetStateVariables()

	if len(stateVariables) == 0 {
		return false
	}

	totalBits := getTotalStorageBits(stateVariables)
	potentialSlots := math.Ceil(float64(totalBits) / maxSlotSize)
	currentSlots := getSlotsUsed(stateVariables)

	if potentialSlots < float64(currentSlots) {
		return true
	}
	return false
}

func getTotalStorageBits(stateVariables []*ir.StateVariable) int64 {
	var sum int64 = 0
	for _, v := range stateVariables {
		storageSize, _ := v.GetStorageSize()
		sum += storageSize
	}
	return sum
}

func getSlotsUsed(stateVariables []*ir.StateVariable) int64 {
	var slotsUsed, bitsUsed int64
	for _, v := range stateVariables {
		size, _ := v.GetStorageSize()
		if size+bitsUsed > maxSlotSize {
			slotsUsed++
			bitsUsed = size
		} else {
			bitsUsed += size
		}
	}
	return slotsUsed + 1
}

func listAllStateVariables(stateVariables []*ir.StateVariable) {
	for _, v := range stateVariables {
		name := v.GetName()
		vartype := v.GetType()
		storageSize, exact := v.GetStorageSize()
		fmt.Println(name, vartype, storageSize, exact)
	}
}
