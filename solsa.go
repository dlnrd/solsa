package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/dlnrd/solgo/printer/ast_printer"
	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/ir"

	opt "solsa/optimisations"
)

func getSources() (sources *solgo.Sources, ok bool) {
	tmpFilePath := flag.String("i", "", "Path to File/Directory")
	outputPath := flag.String("o", "", "Path to save output")
	flag.Parse()

	filePath, _ := fp.Abs(*tmpFilePath)

	var isFile bool

	// CHECKS PATH ISNT EMPTY
	if filePath == "" {
		fmt.Println("Enter a path to file or directory (-i)")
		return nil, false // better than os.Exit(0) as runs deconstructors that close stuff
	}

	// CHECK INPUT PATH VALIDITY
	if fileInfo, err := os.Stat(filePath); err == nil { // file exists
		isFile = !fileInfo.IsDir()
	} else {
		fmt.Println("File doesn't exist, please enter a valid filepath")
		return nil, false
	}
	if *outputPath == "" { // don't want output
		// don't do anything :)
	} else if fileInfo, err := os.Stat(*outputPath); err == nil { // file exists
		if fileInfo.IsDir() {
			fmt.Println("Output file is a directory, please enter a filepath to save")
			return nil, false
		}
	} else { // file doesn't exist
		data := []byte("data")
		if err := os.WriteFile(*outputPath, data, 0644); err == nil {
			os.Remove(*outputPath) // can write to outputfile
		} else {
			fmt.Println("Can't write to output file, do you have write permissions?")
			return nil, false
		}
	}

	if isFile {
		if fp.Ext(filePath) != ".sol" {
			fmt.Println("File is not a .sol file")
			return nil, false
		}

		contractName := strings.TrimSuffix(fp.Base(filePath), ".sol")
		content, err := os.ReadFile(fp.Clean(filePath))
		if err != nil {
			fmt.Println(err)
			return nil, false
		}
		sources = &solgo.Sources{
			SourceUnits: []*solgo.SourceUnit{
				{
					Name:    contractName,
					Path:    filePath,
					Content: string(content),
				},
			},
		}
	} else {
		var err error
		solgo.SetLocalSourcesPath(filePath)
		sources, err = solgo.NewSourcesFromPath("", filePath)
		if err != nil {
			fmt.Println(err)
			return nil, false
		}
	}

	return sources, true
}

func setupSolgoBuilder(sources *solgo.Sources) (builder *ir.Builder, ok bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	builder, err := ir.NewBuilderFromSources(ctx, sources)
	if err != nil {
		fmt.Println("builder init: ", err)
		return nil, false
	}

	if err := builder.Parse(); err != nil {
		fmt.Println("builder parse: ", err)
		return nil, false
	}

	if err := builder.Build(); err != nil {
		fmt.Println("builder build: ", err)
		return nil, false
	}

	ast := builder.GetAstBuilder()
	if err := ast.ResolveReferences(); err != nil {
		fmt.Println("AST Resolve References: ", err)
		return nil, false
	}

	return builder, true
}

func optimiseContracts(builder *ir.Builder) {
	contracts := builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		fmt.Println("\nContract: ", contract.GetName())

		stateVarOpt := opt.StateVariableOptimisable(contract)
		structVarOpt := opt.StructVariableOptimisable(contract)
		calldataOpt := opt.CalldataOptimisable(contract)

		fmt.Println("------------------ OPTIMISATIONS -------------------")
		fmt.Println("StateVariableOptimisable: ", stateVarOpt)
		fmt.Println("StructVariableOptimisable: ", structVarOpt)
		fmt.Println("CalldataOptimisable: ", calldataOpt)

		if stateVarOpt == false && structVarOpt == false && calldataOpt == false { // kinda ugly, maybe fix?
			fmt.Println("-------------- NO OPTIMISATIONS FOUND --------------")
			continue
		}

		unoptContract, _ := ast_printer.Print(contract.GetAST().GetContract())
		fmt.Println("--------------- UNOPTIMISED CONTRACT ---------------")
		fmt.Print(unoptContract)

		if stateVarOpt {
			opt.OptimiseStateVariables(contract)
		}
		if structVarOpt {
			opt.OptimiseStructVariables(contract)
		}
		if calldataOpt {
			opt.OptimiseCalldata(contract)

		}

		optContract, _ := ast_printer.Print(contract.GetAST().GetContract())
		fmt.Println("---------------- OPTIMISED CONTRACT ----------------")
		fmt.Print(optContract)
	}
}

func main() {
	sources, ok := getSources()
	if !ok {
		return
	}

	builder, ok := setupSolgoBuilder(sources)
	if !ok {
		return
	}

	optimiseContracts(builder)

}
