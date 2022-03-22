package main

import (
	"encoding/json"
	"go/ast"
	"go/types"
	"os"

	noglobals "4d63.com/gochecknoglobals/checknoglobals"
	ruleguard "github.com/quasilyte/go-ruleguard/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

const Config = `config/config.json`

type ConfigData struct {
	Staticcheck []string
}

func main() {
	data, err := os.ReadFile(Config)
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	mychecks := []*analysis.Analyzer{
		// ErrCheckAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		unreachable.Analyzer,
		ruleguard.Analyzer,
		noglobals.Analyzer(),
	}

	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Name] || checks["SA*"] {
			mychecks = append(mychecks, v)
		}
	}

	for _, v := range stylecheck.Analyzers {
		if checks[v.Name] || checks["ST*"] {
			mychecks = append(mychecks, v)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}

var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  run,
}

var errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isErrorType(t types.Type) bool {
	return types.Implements(t, errorType)
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isReturnError(pass, call) {
				pass.Reportf(x.Pos(), "expression returns unchecked error")
			}
		}
	}
	tuplefunc := func(x *ast.AssignStmt) {
		if call, ok := x.Rhs[0].(*ast.CallExpr); ok {
			results := resultErrors(pass, call)
			for i := 0; i < len(x.Lhs); i++ {
				if id, ok := x.Lhs[i].(*ast.Ident); ok && id.Name == "_" && results[i] {
					pass.Reportf(id.NamePos, "assignment with unchecked error")
				}
			}
		}
	}
	errfunc := func(x *ast.AssignStmt) {
		for i := 0; i < len(x.Lhs); i++ {
			if id, ok := x.Lhs[i].(*ast.Ident); ok {
				if call, ok := x.Rhs[i].(*ast.CallExpr); ok {
					if id.Name == "_" && isReturnError(pass, call) {
						pass.Reportf(id.NamePos, "assignment with unchecked error")
					}
				}
			}
		}
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.ExprStmt:
				expr(x)
			case *ast.GoStmt:
				if isReturnError(pass, x.Call) {
					pass.Reportf(x.Pos(), "go statement with unchecked error")
				}
			case *ast.DeferStmt:
				if isReturnError(pass, x.Call) {
					pass.Reportf(x.Pos(), "defer with unchecked error")
				}
			case *ast.AssignStmt:
				if len(x.Rhs) == 1 {
					tuplefunc(x)
				} else {
					errfunc(x)
				}
			}
			return true
		})
	}
	return nil, nil
}

func resultErrors(pass *analysis.Pass, call *ast.CallExpr) []bool {
	switch t := pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named:
		return []bool{isErrorType(t)}
	case *types.Pointer:
		return []bool{isErrorType(t)}
	case *types.Tuple:
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch mt := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isErrorType(mt)
			case *types.Pointer:
				s[i] = isErrorType(mt)
			}
		}
		return s
	}
	return []bool{false}
}

func isReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isError := range resultErrors(pass, call) {
		if isError {
			return true
		}
	}
	return false
}
