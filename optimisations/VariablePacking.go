package optimisations

import "sort"

const maxSlotSize = 256 // in bits

type Variable struct {
	Index int
	Size  int64
}

func sumOfBin(bin []Variable) int64 {
	var total int64 = 0
	for _, v := range bin {
		total += v.Size
	}
	return total
}

func VariablePacking(variables []Variable) [][]Variable {
	sort.Slice(variables, func(i, j int) bool { return variables[i].Size > variables[j].Size })
	bins := [][]Variable{}
	for _, v := range variables {
		done := false
		for index, bin := range bins {
			if sumOfBin(bin)+v.Size > maxSlotSize {
				continue
			}
			bins[index] = append(bin, v)
			done = true
			break
		}
		if !done {
			bins = append(bins, []Variable{v})
		}
	}
	return bins
}

// simple bin packing algo
// func ffBinPacking(slice []int) [][]int {
// 	const binsize = 32
// 	bins := [][]int{}
// 	for _, value := range slice {
// 		done := false
// 		for index, binValue := range bins {
// 			if sumOfSlice(binValue)+value > binsize {
// 				continue
// 			}
// 			bins[index] = append(binValue, value)
// 			done = true
// 			break
// 		}
// 		if !done {
// 			bins = append(bins, []int{value})
// 		}
// 	}
// 	return bins
// }
