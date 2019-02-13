package main

import (
	"github.com/tenntenn/gosa/passes/nilerr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(nilerr.Analyzer) }
