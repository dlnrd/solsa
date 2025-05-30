package optimisations

import (
	"strings"

	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/ir"
)

func OptimiseCalldata(contract *ir.Contract) bool {
	functions := contract.GetFunctions()
	if !CalldataOptimisable(contract) {
		return false
	}

	for _, function := range functions {
		if function.GetVisibility() != ast_pb.Visibility_EXTERNAL {
			continue
		}
		params := function.GetAST().GetParameters().GetParameters()

		// if pure, all params can't be modified so should be in calldata
		if function.GetStateMutability() == ast_pb.Mutability_PURE {
			for _, p := range params {
				if p.GetStorageLocation() != ast_pb.StorageLocation_CALLDATA {
					if strings.Contains(p.GetTypeName().GetName(), "[]") {
						p.StorageLocation = ast_pb.StorageLocation_CALLDATA
					}
				}
			}
			continue
		}

		// check for params that need to be optimised
		for _, p := range params {
			if p.GetStorageLocation() != ast_pb.StorageLocation_MEMORY {
				continue
			}

			if !strings.Contains(p.GetTypeName().GetName(), "[]") {
				continue
			}

			if !parameterGetsModified(p, function) {
				p.StorageLocation = ast_pb.StorageLocation_CALLDATA
			}
		}
	}

	return true
}

func CalldataOptimisable(contract *ir.Contract) bool {
	functions := contract.GetFunctions()
	for _, function := range functions {
		if function.GetVisibility() != ast_pb.Visibility_EXTERNAL {
			continue
		}
		params := function.GetAST().GetParameters().GetParameters()

		// if pure, all params can't be modified so should be in calldata
		if function.GetStateMutability() == ast_pb.Mutability_PURE {
			for _, p := range params {
				if p.GetStorageLocation() != ast_pb.StorageLocation_CALLDATA {
					if strings.Contains(p.GetTypeName().GetName(), "[]") {
						return true
					}
				}
			}
		}

		// check for params that need to be optimised
		for _, p := range params {
			if p.GetStorageLocation() != ast_pb.StorageLocation_MEMORY {
				continue
			}
			if !strings.Contains(p.GetTypeName().GetName(), "[]") {
				continue
			}

			if !parameterGetsModified(p, function) {
				return true
			}
		}
	}

	return false
}

func parameterGetsModified(param *ast.Parameter, function *ir.Function) bool {
	assignmentNodes := getAssignmentNodes(function)
	for _, an := range assignmentNodes {
		node := an.(*ast.Assignment)
		nodeName := strings.Split(node.Text, "=")[0]

		if strings.Contains(nodeName, "[") {
			nodeName = nodeName[:strings.Index(nodeName, "[")] // remove index details from slices
		}

		if nodeName == param.Name {
			return true
		}
	}
	return false
}

func getAssignmentNodes(function *ir.Function) []ast.Node[ast.NodeType] {
	assignments := []ast.Node[ast.NodeType]{}
	for _, node := range function.GetAST().GetNodes() {
		assignments = append(assignments, getAssignmentNodesRecursive(node)...)
	}
	return assignments
}

func getAssignmentNodesRecursive(node ast.Node[ast.NodeType]) []ast.Node[ast.NodeType] {
	assignments := []ast.Node[ast.NodeType]{}
	if node.GetNodes() != nil {
		for _, n := range node.GetNodes() {
			assignments = append(assignments, getAssignmentNodesRecursive(n)...)
		}
	}

	if node.GetType() == ast_pb.NodeType_ASSIGNMENT {
		return []ast.Node[ast.NodeType]{node.(*ast.Assignment)}
	}
	return assignments
}
