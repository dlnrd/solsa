package optimisations

import (
	"fmt"
	"math"

	"github.com/unpackdev/solgo/ir"
)

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

func PrintStateVariables(stateVariables []*ir.StateVariable) {
	for _, v := range stateVariables {
		name := v.GetName()
		vartype := v.GetType()
		storageSize, exact := v.GetStorageSize()
		fmt.Println(name, vartype, storageSize, exact)
	}
}
