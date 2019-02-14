package main

import (
	"github.com/tenntenn/gosa/passes/wraperrfmt"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(wraperrfmt.Analyzer) }
