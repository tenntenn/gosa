package wraperrfmt

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ssa"
)

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

var Analyzer = &analysis.Analyzer{
	Name: "wraperrfmt",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

const Doc = "wraperrfmt checks invalid arguments of xerrors.Errorf"

func run(pass *analysis.Pass) (interface{}, error) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for i := range funcs {
		for _, b := range funcs[i].Blocks {
			for _, inst := range b.Instrs {
				if isInvalidErrorf(pass, inst) {
					pass.Reportf(inst.Pos(), "unexpected format. format must end with ': %%w'")
				}
			}
		}
	}

	return nil, nil
}

func isInvalidErrorf(pass *analysis.Pass, inst ssa.Instruction) bool {
	call, ok := inst.(*ssa.Call)
	if !ok {
		return false
	}

	if !isCallErrorf(call) {
		return false
	}

	format := getFormat(call.Call.Args)

	if !strings.Contains(format, "%w") {
		return false
	}

	typ, ok := lastErr(pass, call.Pos())
	if !ok {
		return false
	}

	if strings.HasSuffix(format, ": %w") && typ != nil && types.Implements(typ, errType) {
		return false
	}

	return true
}

func isCallErrorf(call *ssa.Call) bool {

	f := call.Common().StaticCallee()
	if f == nil {
		return false
	}

	if removeVendor(f.Pkg.Pkg.Path()) != "golang.org/x/xerrors" {
		return false
	}

	if f.Name() != "Errorf" {
		return false
	}

	return true
}

func removeVendor(path string) string {
	s := strings.Split(path, "/")
	for i := range s {
		if s[i] == "vendor" {
			return strings.Join(s[i+1:], "/")
		}
	}
	return path
}

func getFormat(args []ssa.Value) string {
	if len(args) == 0 {
		return ""
	}

	format, isConst := args[0].(*ssa.Const)
	if !isConst {
		return ""
	}

	if format.Value.Kind() != constant.String {
		return ""
	}

	return constant.StringVal(format.Value)
}

func lastErr(pass *analysis.Pass, pos token.Pos) (types.Type, bool) {
	file := getFile(pass.Files, pos)
	if file == nil {
		return nil, false
	}

	path, exact := astutil.PathEnclosingInterval(file, pos, pos)
	if !exact || len(path) == 0 {
		return nil, false
	}

	callExpr, ok := path[0].(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	if callExpr.Ellipsis != token.NoPos {
		return nil, false
	}

	if len(callExpr.Args) < 2 {
		return nil, true
	}

	last := callExpr.Args[len(callExpr.Args)-1]
	return pass.TypesInfo.TypeOf(last), true
}

func getFile(fs []*ast.File, pos token.Pos) *ast.File {
	for i := range fs {
		if fs[i].Pos() <= pos && pos <= fs[i].End() {
			return fs[i]
		}
	}
	return nil
}
