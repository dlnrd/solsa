package main

import (
	"flag"
	"fmt"
)

func sumOfSlice(slice []int) int {
	total := 0
	for _, value := range slice {
		total += value
	}
	return total
}

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

	fmt.Println(*filePath, *dirPath, *outputPath)

	// fmt.Println(ffBinPacking([]int{16, 4, 2, 22, 8, 2, 32, 1, 6, 8, 4}))
}

/*

 */
