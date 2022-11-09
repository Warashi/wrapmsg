package wrapmsg

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/Warashi/ssautil"
	"golang.org/x/tools/go/ssa"
)

func formatCall(ctx context.Context, call *ssa.Call) ([]string, bool) {
	c, ok := getCallExpr(ctx, call)
	if !ok {
		return nil, false
	}
	return formatCallExpr(c), true
}

func getCallExpr(ctx context.Context, call ssautil.Poser) (*ast.CallExpr, bool) {
	posMap := ssautil.PosMap(ctx)
	for i := call.Pos(); i > 0; i-- {
		stack := posMap[i]
		if len(stack) == 0 {
			break
		}
		for j := range stack {
			node := stack[len(stack)-1-j]
			ident, ok := node.(*ast.CallExpr)
			if ok {
				return ident, true
			}
		}
	}
	return nil, false
}

func formatCallExpr(call *ast.CallExpr) []string {
	switch f := call.Fun.(type) {
	case *ast.SelectorExpr:
		return formatSelectorExpr(f)
	case *ast.Ident:
		return []string{f.Name}
	}
	return nil
}

func formatSelectorExpr(sel *ast.SelectorExpr) []string {
	var ret []string
	switch x := sel.X.(type) {
	case *ast.CallExpr:
		ret = append(formatCallExpr(x), ret...)
	case *ast.Ident:
		ret = append(ret, x.Name)
	case *ast.SelectorExpr:
		ret = append(formatSelectorExpr(x), ret...)
	case *ast.IndexExpr:
		ret = append(formatIndexExpr(x), ret...)
	}
	ret = append(ret, sel.Sel.Name)
	return ret
}

func formatIndexExpr(expr *ast.IndexExpr) []string {
	var ret []string
	switch x := expr.X.(type) {
	case *ast.SelectorExpr:
		ret = append(formatSelectorExpr(x), ret...)
	}

	switch x := expr.Index.(type) {
	case *ast.Ident:
		ret[len(ret)-1] = fmt.Sprintf("%s[%s]", ret[len(ret)-1], x.Name)
	case *ast.BasicLit:
		ret[len(ret)-1] = fmt.Sprintf("%s[%s]", ret[len(ret)-1], x.Value)
	}
	return ret
}
