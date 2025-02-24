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
	"github.com/unpackdev/solgo/ir"

	opt "solsa/optimisations"
)

func ParseFlags() (filePath, outputPath string, ok bool) {
	tmpFilePath := flag.String("i", "", "Path to File/Directory")
	tmpOutputPath := flag.String("o", "", "Path to save output")
	flag.Parse()

	if *tmpFilePath == "" {
		fmt.Println("Enter a path to file or directory (-i)")
		return "", "", false // better than os.Exit(0) as runs deconstructors that close stuff
	} else {
		filePath, _ = fp.Abs(*tmpFilePath)
	}
	if *tmpOutputPath == "" {
		outputPath = ""
	} else {
		outputPath, _ = fp.Abs(*tmpOutputPath)
	}

	if outputPath == "" { // don't want output
		// don't do anything :)
	} else if fileInfo, err := os.Stat(outputPath); err == nil { // file exists
		if fileInfo.IsDir() {
			fmt.Println("Output file is a directory, please enter a filepath to save")
			return "", "", false
		}
	} else { // file doesn't exist
		data := []byte("data")
		if err := os.WriteFile(outputPath, data, 0644); err == nil {
			os.Remove(outputPath) // can write to outputfile
		} else {
			fmt.Println("Can't write to output file, do you have write permissions?")
			return "", "", false
		}
	}

	return filePath, outputPath, true
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

func OptimiseContracts(builder *ir.Builder) {
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
