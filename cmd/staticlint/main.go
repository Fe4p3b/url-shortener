package main

import (
	"encoding/json"
	"os"

	noglobals "4d63.com/gochecknoglobals/checknoglobals"
	"github.com/Fe4p3b/url-shortener/cmd/staticlint/noosexit"
	gocritic "github.com/go-critic/go-critic/checkers/analyzer"
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

const Config = `config/staticlint/config.json`

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
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		unreachable.Analyzer,
		noglobals.Analyzer(),
		gocritic.Analyzer,
		noosexit.Analyzer(),
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
