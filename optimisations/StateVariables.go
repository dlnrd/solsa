package optimisations

import (
	"fmt"
	"math"

	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/ir"
)

func OptimiseStateVariables(contract *ir.Contract) bool {
	stateVariables := contract.GetStateVariables()
	if !StateVariableOptimisable(contract) {
		return false
	}

	variables := []Variable{}
	// get index and size of variable
	for index, v := range stateVariables {
		size, _ := v.GetStorageSize()
		variables = append(variables, Variable{index, size})
	}
	packedVariables := VariablePacking(variables)

	// convert 2d array into *[]ir.StateVariable
	newStateVariables := []*ir.StateVariable{}
	for _, bin := range packedVariables {
		for _, v := range bin {
			newStateVariables = append(newStateVariables, stateVariables[v.Index])
		}
	}

	// get parent nodes (make this more readable)
	parentNodes := newStateVariables[0].GetAST().GetTree().GetById(newStateVariables[0].GetSrc().GetParentIndex()).GetNodes()

	tmpNodes := make([]ast.Node[ast.NodeType], len(parentNodes))
	copy(tmpNodes, parentNodes)

	m := make(map[int64]int)

	// init map for sorting
	for destinationIndex, newStateVariable := range newStateVariables {
		childId := newStateVariable.GetAST().GetId()
		m[childId] = destinationIndex
	}

	// rearrange state variable nodes in AST
	for index, child := range tmpNodes {
		destIndex, exist := m[child.GetId()]
		if !exist {
			continue
		}
		parentNodes[destIndex] = tmpNodes[index]

	}
	// update ir with new variables
	contract.StateVariables = newStateVariables
	return true
}

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
