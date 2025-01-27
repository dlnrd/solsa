package optimisations

import (
	"fmt"
	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/ir"
	"math"
)

const maxSlotSize = 256 // in bits

func StructVariableOptimisable(contract *ir.Contract) bool {
	structs := contract.GetStructs()
	if len(structs) == 0 {
		return false
	}

	for _, s := range structs {
		structVariables := s.GetAST().GetMembers()
		getTotalStorageBitsStruct(structVariables)
		totalBits := getTotalStorageBitsStruct(structVariables)
		potentialSlots := math.Ceil(float64(totalBits) / maxSlotSize)
		currentSlots := getSlotsUsedStruct(structVariables)
		if potentialSlots < float64(currentSlots) {
			return true
		}
	}
	return false
}

func getTotalStorageBitsStruct(structVariables []*ast.Parameter) int64 {
	var sum int64 = 0
	for _, v := range structVariables {
		storageSize, _ := v.GetTypeName().StorageSize()
		sum += storageSize
	}
	return sum
}

func getSlotsUsedStruct(structVariables []*ast.Parameter) int64 {
	var slotsUsed, bitsUsed int64
	for _, v := range structVariables {
		size, _ := v.GetTypeName().StorageSize()
		if size+bitsUsed > maxSlotSize {
			slotsUsed++
			bitsUsed = size
		} else {
			bitsUsed += size
		}
	}
	return slotsUsed + 1
}

func PrintStructVariables(contract *ir.Contract) {
	structs := contract.GetStructs()
	for _, s := range structs {
		members := s.GetAST().GetMembers()
		for _, param := range members {
			name := param.GetName()
			vartype := param.GetTypeName().GetName()
			size, exact := param.GetTypeName().StorageSize()
			fmt.Println(name, vartype, size, exact)
		}
	}
}
