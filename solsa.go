package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/ir"

	opt "solsa/optimisations"
)

/*
 TODO:
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

	var sources *solgo.Sources
	if isFile {
		if fp.Ext(filePath) != ".sol" {
			fmt.Println("File is not a .sol file")
			return
		}

		contractName := strings.TrimSuffix(fp.Base(filePath), ".sol")
		content, err := os.ReadFile(fp.Clean(filePath))
		if err != nil {
			fmt.Println(err)
			return
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
			return
		}
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

	contracts := builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		fmt.Println("Contract: ", contract.GetName())
		fmt.Println("StateVariableOptimisable: ", opt.StateVariableOptimisable(contract))
		fmt.Println("StructVariableOptimisable: ", opt.StructVariableOptimisable(contract))
	}
}

// var declarations struct {
// 	Name string
// 	VarType string
// 	Size int64
// 	Exact bool
// }
