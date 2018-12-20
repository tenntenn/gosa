package findfuncbody

import (
	"reflect"

	"github.com/tenntenn/gosa/passes/findfuncbody/funcbody"
	"github.com/tenntenn/gosa/passes/internal/buildssa"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:       "findfuncbody",
	Doc:        Doc,
	Run:        run,
	ResultType: reflect.TypeOf(new(funcbody.Finder)),
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

const Doc = "findfuncbody is ..."

func run(pass *analysis.Pass) (interface{}, error) {
	ssa := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).Pkg
	return &funcbody.Finder{
		Fset:      pass.Fset,
		Files:     pass.Files,
		TypesInfo: pass.TypesInfo,
		Pkg:       pass.Pkg,
		SSA:       ssa,
	}, nil
}
