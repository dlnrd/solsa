package optimisations

import (
	"fmt"
	"math"

	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/ir"
)

func OptimiseStateVariables(contract *ir.Contract) bool {
	stateVariables := contract.GetStateVariables()
	if len(stateVariables) == 0 {
		return false
	}

	variables := []Variable{}
	// get index and size of variable
	for index, v := range stateVariables {
		size, _ := v.GetStorageSize()
		variables = append(variables, Variable{index, size})
	}
	tmp := VariablePacking(variables)
	// convert 2d array into *[]ir.StateVariable
	newStateVariables := []*ir.StateVariable{}
	for _, bin := range tmp {
		for _, v := range bin {
			newStateVariables = append(newStateVariables, stateVariables[v.Index])
		}
	}

	tree := newStateVariables[0].GetAST().GetTree()
	parent_node := tree.GetById(newStateVariables[0].GetSrc().GetParentIndex())
	parent_nodes := parent_node.GetNodes()

	cp := make([]ast.Node[ast.NodeType], len(parent_nodes))
	copy(cp, parent_nodes)

	m := make(map[int64]int)

	for desindex, nsv := range newStateVariables {
		child_id := nsv.GetAST().GetId()
		m[child_id] = desindex
	}

	for index, child := range cp {
		destIndex, exist := m[child.GetId()]
		if !exist {
			continue
		}
		parent_nodes[destIndex] = cp[index]

	}
	// update ast with new variables
	contract.StateVariables = newStateVariables
	return true
}

func printAllNodes(node ast.Node[ast.NodeType]) {
	nodes := node.GetNodes()
	if nodes != nil {
		for i, node := range nodes {
			fmt.Println(i)
			printAllNodes(node)
		}
	}
	fmt.Println(node)
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
