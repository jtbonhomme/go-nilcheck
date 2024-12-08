package main

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// Analyzer is the linter's main logic.
var Analyzer = &analysis.Analyzer{
	Name: "nilcheck",
	Doc:  "checks if functions with pointer arguments test for nil",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass == nil {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Look for function declarations
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Check if the function has pointer arguments
			pointerArgs := getPointerArguments(pass, funcDecl)
			if len(pointerArgs) == 0 {
				return true
			}

			// Check for nil checks within the function body
			nilChecked := checkNilInBody(funcDecl.Body, pointerArgs)
			if !nilChecked {
				pass.Reportf(funcDecl.Pos(), "Function %s has pointer arguments but does not check for nil", funcDecl.Name.Name)
			}

			return true
		})
	}
	return nil, nil
}

// getPointerArguments extracts pointer arguments from a function.
func getPointerArguments(pass *analysis.Pass, funcDecl *ast.FuncDecl) map[string]struct{} {
	pointerArgs := make(map[string]struct{})
	if pass == nil {
		return pointerArgs
	}
	if funcDecl == nil {
		return pointerArgs
	}
	if funcDecl.Type.Params == nil {
		return pointerArgs
	}

	for _, param := range funcDecl.Type.Params.List {
		paramType := pass.TypesInfo.TypeOf(param.Type)
		if _, ok := paramType.(*types.Pointer); ok {
			// Add each pointer argument name to the map
			for _, name := range param.Names {
				pointerArgs[name.Name] = struct{}{}
			}
		}
	}
	return pointerArgs
}

// checkNilInBody checks if the function body tests pointer arguments for nil.
func checkNilInBody(body *ast.BlockStmt, pointerArgs map[string]struct{}) bool {
	if body == nil {
		return false
	}

	foundNilCheck := false
	ast.Inspect(body, func(n ast.Node) bool {
		if foundNilCheck {
			return false // Stop further inspection once a nil check is found
		}

		// Look for if statements
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		// Check if the condition is a nil check for any pointer argument
		binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		if binaryExpr.Op == token.EQL || binaryExpr.Op == token.NEQ {
			ident, ok := binaryExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			if _, isPointerArg := pointerArgs[ident.Name]; isPointerArg {
				// Check if the comparison involves "nil"
				if isNil(binaryExpr.Y) || isNil(binaryExpr.X) {
					foundNilCheck = true
				}
			}
		}
		return true
	})
	return foundNilCheck
}

// isNil checks if the expression is a nil literal.
func isNil(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

func main() {
	singlechecker.Main(Analyzer)
}
