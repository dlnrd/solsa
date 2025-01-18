package main

import (
	"flag"
	"fmt"
	"os"
	// "path/filepath"
)

func sumOfSlice(slice []int) int {
	total := 0
	for _, value := range slice {
		total += value
	}
	return total
}

// simple bin packing algo
func ffBinPacking(slice []int) [][]int {
	const binsize = 32
	bins := [][]int{}
	for _, value := range slice {
		done := false
		for index, binValue := range bins {
			if sumOfSlice(binValue)+value > binsize {
				continue
			}
			bins[index] = append(binValue, value)
			done = true
			break
		}
		if !done {
			bins = append(bins, []int{value})
		}
	}
	return bins
}

func main() {
	filePath := flag.String("f", "", "Path to File")
	dirPath := flag.String("d", "", "Path to Directory")
	outputPath := flag.String("o", "", "Path to save output")
	flag.Parse()

	if *filePath == "" && *dirPath == "" {
		fmt.Println("Please enter a path to either a file (-f) or directory (-d)")
		return
	}
	if *filePath != "" && *dirPath != "" {
		fmt.Println("Please use either -f or -d not both")
		return // better than os.Exit(0) as runs deconstructors that close stuff
	}
	// if !filepath.IsLocal(*filePath) && !filepath.IsAbs(*filePath) {
	// 	return // not sure this is best to verify input is filepath
	// }
	if i, err := os.Stat(*filePath); err == nil {
		fmt.Println("File exists")
		fmt.Println(i.IsDir())

	}

	fmt.Println("FilePath: ", *filePath, "\tDirPath: ", *dirPath, "\tOutputPath: ", *outputPath)

	// fmt.Println(ffBinPacking([]int{16, 4, 2, 22, 8, 2, 32, 1, 6, 8, 4}))
}

/*

 */
