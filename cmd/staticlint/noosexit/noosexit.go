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
			case *ast.FuncDecl:
				if x.Name.String() == "main" {
					filterBlockStmt(pass, x.Body)
				}
			}

			return true
		})
	}
	return nil, nil
}

// filterBlockStmt traverses ast.BlockStmt for if statement
// or expresstion.
func filterBlockStmt(pass *analysis.Pass, body *ast.BlockStmt) {
	for _, node := range body.List {
		switch x := node.(type) {
		case *ast.ExprStmt:
			nodeWithExit(pass, x)
		case *ast.IfStmt:
			filterBlockStmt(pass, x.Body)
		}
	}
}

// nodeWithExit finds os.Exit call in node and
// reports about it.
func nodeWithExit(pass *analysis.Pass, node *ast.ExprStmt) {
	switch x := node.X.(type) {
	case *ast.CallExpr:
		if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
			if pkg, ok := fun.X.(*ast.Ident); ok {
				funcName := fun.Sel.Name
				if pkg.String() == "os" && funcName == "Exit" {
					pass.Reportf(fun.X.Pos(), "os.Exit in main")
				}
			}

		}
	}
}
