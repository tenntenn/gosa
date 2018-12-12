package findfuncbody

import (
	"reflect"

	"github.com/tenntenn/gosa/passes/findfuncbody/funcbody"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:       "findfuncbody",
	Doc:        Doc,
	Run:        run,
	ResultType: reflect.TypeOf(new(funcbody.Finder)),
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "findfuncbody is ..."

func run(pass *analysis.Pass) (interface{}, error) {
	return &funcbody.Finder{
		TypesInfo: pass.TypesInfo,
		Fset:      pass.Fset,
	}, nil
}
