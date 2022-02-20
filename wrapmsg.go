package wrapmsg

import (
	"context"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/ssa"
)

const doc = "wrapmsg is linter for error-wrapping message"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "wrapmsg",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		buildssa.Analyzer,
	},
}

type poser interface {
	Pos() token.Pos
}

type posReferrerer interface {
	poser
	Referrers() *[]ssa.Instruction
}

type posOperander interface {
	poser
	operander
}

func isErrorf(call *ssa.Call) bool {
	if f, ok := getOperands(call)[0].(*ssa.Function); ok && f.Pkg.Pkg.Path() == "testing" {
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
	for _, f := range getSSA(ctx).SrcFuncs {
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

func prepare(ctx context.Context, pass *analysis.Pass) context.Context {
	ctx = context.WithValue(ctx, passKey, pass)
	ctx = context.WithValue(ctx, ssaKey, pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA))
	ctx = context.WithValue(ctx, inspectorKey, pass.ResultOf[inspect.Analyzer].(*inspector.Inspector))
	ctx = context.WithValue(ctx, posMapKey, buildPosMap(ctx))
	return ctx
}

func buildPosMap(ctx context.Context) map[token.Pos][]ast.Node {
	posMap := make(map[token.Pos][]ast.Node)
	getInspector(ctx).Preorder(nil, func(node ast.Node) {
		for i := node.Pos(); i <= node.End(); i++ {
			posMap[i] = append(posMap[i], node)
		}
	})
	return posMap
}

func run(pass *analysis.Pass) (interface{}, error) {
	ctx := prepare(context.Background(), pass)
	for _, call := range iterateErrorf(ctx) {
		report(ctx, call)
	}
	return nil, nil
}
