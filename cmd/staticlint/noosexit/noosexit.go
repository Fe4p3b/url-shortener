// Package noosexit implements analizer with
// basic functionality of finding os.Exit
// function call in main function.
package noosexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer creates analysis.Analyzer
func Analyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "noosexit",
		Doc:  "check for os.Exit in main func",
		Run:  run,
	}
}

// run traverses files with filename main.go, finds
// main function with os.Exit call. Then displays
// position and file of os.exit function.
func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		if !strings.HasSuffix(filename, "main.go") {
			continue
		}

		if strings.HasSuffix(filename, "_test.go") {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.ExprStmt:
				nodeWithExit(x, pass)
			}

			return true
		})
	}
	return nil, nil
}

// nodeWithExit finds os.Exit call in node and
// reports about it.
func nodeWithExit(node ast.Stmt, pass *analysis.Pass) {
	if node, ok := node.(*ast.ExprStmt); ok {
		if call, ok := node.X.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := fun.X.(*ast.Ident); ok {
					funcName := fun.Sel.Name
					if pkg.String() == "os" && funcName == "Exit" {
						pass.Reportf(fun.X.Pos(), "os.Exit in main")
					}
				}

			}
		}
	}
}
