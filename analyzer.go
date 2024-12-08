package nilcheck

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the linter's main logic.
var Analyzer = &analysis.Analyzer{
	Name: "nilpointerlinter",
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
	if pass == nil ||
		funcDecl == nil ||
		funcDecl.Type.Params == nil {
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

	// Track which pointer arguments have been tested for nil
	testedPointers := make(map[string]struct{})
	ast.Inspect(body, func(n ast.Node) bool {
		// Look for if statements
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		// Check the condition for nil checks
		checkNilInCondition(ifStmt.Cond, pointerArgs, testedPointers)
		return true
	})

	// Verify all pointer arguments are tested
	for arg := range pointerArgs {
		if _, tested := testedPointers[arg]; !tested {
			return false
		}
	}

	return true
}

// checkNilInCondition recursively checks conditions for nil checks and tracks tested arguments.
func checkNilInCondition(expr ast.Expr, pointerArgs map[string]struct{}, testedPointers map[string]struct{}) {
	switch v := expr.(type) {
	case *ast.BinaryExpr:
		// Handle binary expressions
		switch v.Op {
		case token.EQL, token.NEQ:
			// Check for equality or inequality with nil
			if ident, ok := v.X.(*ast.Ident); ok {
				if _, isPointerArg := pointerArgs[ident.Name]; isPointerArg && isNil(v.Y) {
					testedPointers[ident.Name] = struct{}{}
				}
			}
			if ident, ok := v.Y.(*ast.Ident); ok {
				if _, isPointerArg := pointerArgs[ident.Name]; isPointerArg && isNil(v.X) {
					testedPointers[ident.Name] = struct{}{}
				}
			}
		case token.LOR, token.LAND:
			// Handle logical OR/AND (e.g., p == nil || q == nil)
			checkNilInCondition(v.X, pointerArgs, testedPointers)
			checkNilInCondition(v.Y, pointerArgs, testedPointers)
		}
	case *ast.ParenExpr:
		// Handle parenthesized expressions (e.g., (p == nil))
		checkNilInCondition(v.X, pointerArgs, testedPointers)
	}
}

// isNil checks if the expression is a nil literal.
func isNil(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}
