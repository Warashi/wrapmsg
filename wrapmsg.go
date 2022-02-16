package wrapmsg

import (
	"context"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"
	"strings"

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

func replaceConst(expr *ast.CallExpr, actual, want string) *ast.CallExpr {
	ret := *expr
	ret.Args = make([]ast.Expr, len(expr.Args))
	copy(ret.Args, expr.Args)

	for i, arg := range ret.Args {
		c, ok := arg.(*ast.BasicLit)
		if !ok {
			continue
		}
		if c.Kind != token.STRING {
			continue
		}
		if strconv.Quote(actual) == c.Value || (strconv.CanBackquote(actual) && c.Value == "`"+actual+"`") {
			ret.Args[i] = &ast.BasicLit{
				ValuePos: c.ValuePos,
				Kind:     c.Kind,
				Value:    strconv.Quote(want),
			}
			return &ret
		}
	}

	return &ret
}

func genText(ctx context.Context, expr *ast.CallExpr) []byte {
	return []byte(prettyPrint(ctx, expr))
}

func report(ctx context.Context, call *ssa.Call) {
	pass := getPass(ctx)
	var actual, want string
	var gotActual, gotWant bool
	for _, v := range getOperands(call) {
		switch v := v.(type) {
		case *ssa.Const:
			if v.Value == nil || v.Value.Kind() != constant.String {
				continue
			}
			val := constant.StringVal(v.Value)
			if !strings.Contains(val, "%w") {
				continue
			}
			if !gotActual {
				actual = val
				gotActual = true
			}
		case *ssa.Slice:
			w := new(walker)
			if r, ok := w.walk(ctx, 0, v); ok && len(r) > 0 {
				want = strings.Join(r, ".") + ": %w"
				gotWant = true
			}
		}
	}
	if gotWant && gotActual && actual != want {
		node, ok := getCallExpr(ctx, call)
		if !ok {
			return
		}
		pos, end := node.Pos(), node.End()
		pass.Report(analysis.Diagnostic{
			Pos:     pos,
			End:     end,
			Message: fmt.Sprintf("want `the error-wrapping message should be %q", want),
			SuggestedFixes: []analysis.SuggestedFix{{TextEdits: []analysis.TextEdit{{
				Pos:     pos,
				End:     end,
				NewText: genText(ctx, replaceConst(node, actual, want)),
			}}}},
		})
	}
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
	ctx := context.Background()
	ctx = context.WithValue(ctx, passKey, pass)
	ctx = context.WithValue(ctx, ssaKey, pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA))
	ctx = context.WithValue(ctx, inspectorKey, pass.ResultOf[inspect.Analyzer].(*inspector.Inspector))
	ctx = context.WithValue(ctx, posMapKey, buildPosMap(ctx))
	for _, call := range iterateErrorf(ctx) {
		report(ctx, call)
	}
	return nil, nil
}
