package solsa

import (
	"context"
	"flag"
	"fmt"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/dlnrd/solgo/printer/ast_printer"
	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/ir"

	opt "solsa/optimisations"
)

func ParseFlags() (filePath, outputPath string, silent, ok bool) {
	flag.StringVar(&filePath, "i", "", "Path to File/Directory")
	flag.StringVar(&outputPath, "o", "", "Path to output directory")
	flag.BoolVar(&silent, "s", false, "Silent mode, no output to stdout")
	flag.Parse()

	if filePath == "" {
		fmt.Println("Enter a path to file or directory (-i)")
		return "", "", false, false // better than os.Exit(0) as runs deconstructors that close stuff
	} else {
		filePath, _ = fp.Abs(filePath)
	}
	if outputPath != "" {
		outputPath, _ = fp.Abs(outputPath)
	}

	if outputPath == "" {
		// don't want output, don't do anything :)
	} else if fileInfo, err := os.Stat(outputPath); err == nil { // file exists
		if !fileInfo.IsDir() {
			fmt.Println("Output file is a directory, please enter a filepath to save")
			return "", "", false, false
		}
	} else { // file doesn't exist
		data := []byte("data")
		if err := os.WriteFile(outputPath, data, 0644); err == nil {
			os.Remove(outputPath) // can write to outputfile
		} else {
			fmt.Println("Can't write to output file, do you have write permissions?")
			return "", "", false, false
		}
	}

	return filePath, outputPath, silent, true
}

func GetSources(filePath string) (sources *solgo.Sources, ok bool) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("File doesn't exist, please enter a valid filepath")
		return nil, false
	}

	if fileInfo.IsDir() {
		solgo.SetLocalSourcesPath(filePath)
		sources, err = solgo.NewSourcesFromPath("", filePath)
		if err != nil {
			fmt.Println(err)
			return nil, false
		}
	} else {
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
	}

	return sources, true
}

func SetupSolgoBuilder(sources *solgo.Sources) (builder *ir.Builder, ok bool) {
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

func OptimiseContracts(builder *ir.Builder, silent bool) (oState, rState, oStruct, rStruct, oCall, rCall, unopt, failed, totalContracts int) {
	contracts := builder.GetRoot().GetContracts()
	totalContracts = len(contracts)
	for _, contract := range contracts {

		stateVarOpt := opt.StateVariableOptimisable(contract)
		structVarOpt := opt.StructVariableOptimisable(contract)
		calldataOpt := opt.CalldataOptimisable(contract)

		if !silent {
			fmt.Println("\nContract: ", contract.GetName())
			fmt.Println("------------------ OPTIMISATIONS -------------------")
			fmt.Println("StateVariableOptimisable: ", stateVarOpt)
			fmt.Println("StructVariableOptimisable: ", structVarOpt)
			fmt.Println("CalldataOptimisable: ", calldataOpt)
		}

		if stateVarOpt == false && structVarOpt == false && calldataOpt == false { // kinda ugly, maybe fix?
			if !silent {
				fmt.Println("-------------- NO OPTIMISATIONS FOUND --------------")
			}
			unopt++
			continue
		}

		unoptContract, ok := ast_printer.Print(contract.GetAST().GetContract())
		if unoptContract == "" {
			unoptContract = "ERROR: ast_printer cannot print contract\n"
			if !ok { // debug
				// 	fmt.Println("\nContract: ", contract.GetName())
				// 	printAllNodes(contract)
			}
		}

		if stateVarOpt {
			oState++
			opt.OptimiseStateVariables(contract)
		}
		if structVarOpt {
			oStruct++
			opt.OptimiseStructVariables(contract)
		}
		if calldataOpt {
			oCall++
			opt.OptimiseCalldata(contract)

		}

		optContract, _ := ast_printer.Print(contract.GetAST().GetContract())
		if optContract == "" {
			optContract = "ERROR: ast_printer cannot print contract\n"
			failed++
		} else {
			if stateVarOpt {
				rState++
			}
			if structVarOpt {
				rStruct++
			}
			if calldataOpt {
				rCall++

			}
		}

		if !silent {
			fmt.Println("--------------- UNOPTIMISED CONTRACT ---------------")
			fmt.Print(unoptContract)
			fmt.Println("---------------- OPTIMISED CONTRACT ----------------")
			fmt.Print(optContract)
		}
		// fmt.Println(ContractBuilder(contract))
		// builder.Build()
		// builder.GetSources().WriteToDir("/home/dlnrd/uni/fyp/testsol/tmp")

	}

	return oState, rState, oStruct, rStruct, oCall, rCall, unopt, failed, totalContracts
}

func printAllNodes(contract *ir.Contract) {
	contractNodes := contract.GetAST().GetNodes()
	for _, node := range contractNodes {
		printAllNodesRecursive(node)
	}
}

func printAllNodesRecursive(node ast.Node[ast.NodeType]) {
	fmt.Println(node)
	if node.GetNodes() != nil {
		for _, childNode := range node.GetNodes() {
			if childNode != nil {
				printAllNodesRecursive(childNode)
			}
		}
	}
}

func ContractBuilder(contract *ir.Contract) (contractName, optContract string) {
	optContract = "// SPDX-License-Identifier: " + contract.GetLicense() + "\n"
	pragmas := contract.GetPragmas()
	for _, pragma := range pragmas {
		optContract += pragma.GetText() + "\n"
	}
	optContract += "\n"

	body, _ := ast_printer.Print(contract.GetAST().GetContract())
	optContract += body

	return contract.GetName(), optContract
}

// func WriteContracts(outputPath, contractName string, optContracts []string) bool {
// 	if outputPath == "" {
// 		return false
// 	}
// 	contract := []byte(optContracts[])
// 	os.WriteFile(contractName, ,0644)
//
// 	return true
// }
