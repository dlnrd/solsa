package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/unpackdev/solgo"
	"github.com/unpackdev/solgo/ast"
	"os"
	fp "path/filepath"
)

/*
 TODO:
 - parse .sol file
 - create abstract syntax tree (ast)
 - find size of the contracts variables
 - use bin packing algorithm to sort into more efficiently packed slots
 - rewrite bin packing algo to potentially find better solutions (less bins)
 - output results
 - write more optimisations :)
*/

func main() {
	filePath := flag.String("i", "", "Path to File/Directory")
	outputPath := flag.String("o", "", "Path to save output")
	flag.Parse()

	var isFile bool

	// CHECKS PATH ISNT EMPTY
	if *filePath == "" {
		fmt.Println("Enter a path to file or directory (-i)")
		return // better than os.Exit(0) as runs deconstructors that close stuff
	}

	// CHECK INPUT PATH VALIDITY
	if fileInfo, err := os.Stat(*filePath); err == nil { // file exists
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
		if fp.Ext(*filePath) != ".sol" {
			fmt.Println("File is not a .sol file")
			return
		}
		// cwd, _ := os.Getwd()
		// dir := fp.Dir(*filePath)
		// name := fp.Base(*filePath)
		// fmt.Println(cwd, dir, name)
		// fmt.Println(fp.Clean(*filePath))
		// fmt.Println(fp.Join(cwd, name))
	} else {
		// get a list of all the files to be checked?
	}

	// file, err := os.Open(*filePath)
	// if err != nil {
	// 	fmt.Println("Error opening file")
	// 	return
	// }
	// defer file.Close()

	// // boilerplate for solgo
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	cwd, _ := os.Getwd()
	// ffp := fp.Join(cwd, fp.Base(*filePath))

	solgo.SetLocalSourcesPath(cwd)
	sources, err := solgo.NewSourcesFromPath("test", cwd)
	if err != nil {
		fmt.Println(err)
		return
	}

	parser, err := solgo.NewParserFromSources(context.TODO(), sources)
	if err != nil {
		fmt.Println(err)
		return
	}
	astBuilder := ast.NewAstBuilder(parser.GetParser(), parser.GetSources())
	fmt.Println(astBuilder.GetRoot())

	fmt.Println("FilePath: ", *filePath, "\tOutputPath: ", *outputPath)
}
