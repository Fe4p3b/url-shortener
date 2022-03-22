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
			if x, ok := node.(*ast.FuncDecl); ok {
				if x.Name.String() == "main" {
					report(pass, x)
				}
			}

			return true
		})
	}
	return nil, nil
}

// report finds os.Exit call in *ast.FuncDecl and
// reports about it.
func report(pass *analysis.Pass, decl *ast.FuncDecl) {
	for _, node := range decl.Body.List {
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
}
