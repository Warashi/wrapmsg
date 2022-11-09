package wrapmsg

import (
	"context"

	"github.com/Warashi/ssautil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

const doc = "wrapmsg is linter for error-wrapping message"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name:     "wrapmsg",
	Doc:      doc,
	Run:      run,
	Requires: ssautil.Requires(),
}

func isErrorf(call *ssa.Call) bool {
	if f, ok := ssautil.Operands(call)[0].(*ssa.Function); ok && f.Pkg.Pkg.Path() == "testing" {
		// avoid targeting (*testing.T).Errorf
		return false
	}
	if f := call.Common().Method; f != nil {
		return f.Name() == "Errorf"
	}
	if f := call.Common().StaticCallee(); f != nil {
		return f.Name() == "Errorf"
	}
	return false
}

func iterateErrorf(ctx context.Context) []*ssa.Call {
	var r []*ssa.Call
	for _, f := range ssautil.SSA(ctx).SrcFuncs {
		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				switch v := instr.(type) {
				case *ssa.Call:
					if isErrorf(v) {
						r = append(r, v)
					}
				}
			}
		}
	}
	return r
}

func run(pass *analysis.Pass) (interface{}, error) {
	ctx := ssautil.Prepare(context.Background(), pass)
	for _, call := range iterateErrorf(ctx) {
		report(ctx, call)
	}
	return nil, nil
}
