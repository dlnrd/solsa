package optimisations

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/ast"
	"github.com/unpackdev/solgo/ir"
)

func OptimiseCalldata(contract *ir.Contract) {
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
					p.StorageLocation = ast_pb.StorageLocation_CALLDATA
				}
			}
			continue
		}

		// check for params that need to be optimised
		for _, p := range params {
			if p.GetStorageLocation() != ast_pb.StorageLocation_MEMORY {
				continue
			} // param needs to be in memory
			if !strings.Contains(p.GetTypeName().GetName(), "[]") {
				continue
			} // and an array

			if !parameterGetsModified(p, function) {
				p.StorageLocation = ast_pb.StorageLocation_CALLDATA
			}
		}

	}
}

func CalldataOptimisable(contract *ir.Contract) bool {
	functions := contract.GetFunctions()
	if len(functions) == 0 {
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
					return true
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

// TODO: Do this properly, bad code :(
func parameterGetsModified(param *ast.Parameter, function *ir.Function) bool {
	regex := ";" + param.GetName() + "\\[.+\\]="
	matched, err := regexp.MatchString(regex, function.GetAST().ToString())
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return matched
}
