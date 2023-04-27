package wrapmsg

import (
	"context"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"
	"strings"

	"github.com/Warashi/ssautil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

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
	return []byte(ssautil.PrettyPrint(ctx, expr))
}

func report(ctx context.Context, call *ssa.Call) {
	pass := ssautil.Pass(ctx)
	var actual, want string
	var gotActual, gotWant bool
	for _, v := range ssautil.Operands(call) {
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
			if r, ok := w.walk(ctx, v); ok && len(r) > 0 {
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
			Pos:      pos,
			End:      end,
			Category: "wrapmsg",
			Message:  fmt.Sprintf("want `the error-wrapping message should be %q", want),
			SuggestedFixes: []analysis.SuggestedFix{{TextEdits: []analysis.TextEdit{{
				Pos:     pos,
				End:     end,
				NewText: genText(ctx, replaceConst(node, actual, want)),
			}}}},
		})
	}
}
