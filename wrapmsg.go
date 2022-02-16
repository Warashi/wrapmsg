package wrapmsg

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
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

type walker struct {
	stack []interface{}
}

func (w *walker) push(n interface{}) {
	w.stack = append(w.stack, n)
}

func (w *walker) pop() {
	w.stack = w.stack[:len(w.stack)-1]
}

func (w *walker) contains(n interface{}) bool {
	for _, s := range w.stack {
		if s == n {
			return true
		}
	}
	return false
}

func getIdent(ctx context.Context, v poser) (*ast.Ident, bool) {
	posMap := ctx.Value(posMapKey).(map[token.Pos][]ast.Node)
	for _, node := range posMap[v.Pos()] {
		ident, ok := node.(*ast.Ident)
		if ok {
			return ident, true
		}
	}
	return nil, false
}

func getIdentName(ctx context.Context, v poser) []string {
	ident, ok := getIdent(ctx, v)
	switch v := v.(type) {
	case *ssa.Slice:
		return nil
	case *ssa.Alloc:
		switch v.Comment {
		case "varargs":
			return nil
		}
		break
	case *ssa.IndexAddr:
		break
	case *ssa.FieldAddr:
		return nil
	case *ssa.Store:
		return nil
	case *ssa.ChangeInterface:
		return nil
	case *ssa.Call:
		break
	case *ssa.UnOp:
		break
	case *ssa.Parameter:
		return []string{v.Object().Name()}
	case *ssa.Function:
		break
	case *ast.Ident:
		return []string{v.Name}
	case *ast.SelectorExpr:
		return nil
	case *ast.CallExpr:
		return nil
	default:
		return nil
	}

	if !ok {
		return nil
	}

	return []string{ident.Name}
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
	Operander
}

func (w *walker) walkRefs(ctx context.Context, depth int, v posReferrerer) ([]string, bool) {
	org := v
	for _, v := range *v.Referrers() {
		if r, ok := w.walk(ctx, depth, v); ok {
			return append(r, getIdentName(ctx, org)...), true
		}
	}
	return nil, false
}

func (w *walker) walkOperands(ctx context.Context, depth int, v posOperander) ([]string, bool) {
	org := v
	for _, v := range GetOperands(v) {
		if r, ok := w.walk(ctx, depth, v); ok {
			return append(r, getIdentName(ctx, org)...), true
		}
	}
	return nil, false
}

func prettyPrint(ctx context.Context, expr ast.Expr) string {
	pass := ctx.Value(passKey).(*analysis.Pass)
	var b bytes.Buffer
	format.Node(&b, pass.Fset, expr)
	return b.String()
}

func (w *walker) walk(ctx context.Context, depth int, v poser) ([]string, bool) {
	if w.contains(v) {
		return nil, false
	}
	w.push(v)
	defer w.pop()

	switch v := v.(type) {
	case *ssa.Const:
	case *ssa.Slice:
		return w.walkOperands(ctx, depth+1, v)
	case *ssa.Alloc:
		return w.walkRefs(ctx, depth+1, v)
	case *ssa.IndexAddr:
		return w.walkRefs(ctx, depth+1, v)
	case *ssa.Store:
		return w.walkOperands(ctx, depth+1, v)
	case *ssa.ChangeInterface:
		return w.walkOperands(ctx, depth+1, v)
	case *ssa.Call:
		call, ok := getCallExpr(ctx, v)
		if ok {
			return formatCallExpr(call), true
		}
	}
	return nil, false
}

func buildPosMap(ctx context.Context) map[token.Pos][]ast.Node {
	posMap := make(map[token.Pos][]ast.Node)
	i := ctx.Value(inspectorKey).(*inspector.Inspector)
	i.Preorder(nil, func(node ast.Node) {
		for i := node.Pos(); i <= node.End(); i++ {
			posMap[i] = append(posMap[i], node)
		}
	})
	return posMap
}

func isErrorf(call *ssa.Call) bool {
	if f, ok := GetOperands(call)[0].(*ssa.Function); ok && f.Pkg.Pkg.Path() == "testing" {
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
	s := ctx.Value(ssaKey).(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
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

func formatSelectorExpr(sel *ast.SelectorExpr) []string {
	var ret []string
	switch x := sel.X.(type) {
	case *ast.CallExpr:
		ret = append(formatCallExpr(x), ret...)
	case *ast.Ident:
		ret = append(ret, x.Name)
	case *ast.SelectorExpr:
		ret = append(formatSelectorExpr(x), ret...)
	}
	ret = append(ret, sel.Sel.Name)
	return ret
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

func getCallExpr(ctx context.Context, call poser) (*ast.CallExpr, bool) {
	posMap := ctx.Value(posMapKey).(map[token.Pos][]ast.Node)
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

func replaceConst(expr *ast.CallExpr, actual, want string) *ast.CallExpr {
	ret := new(ast.CallExpr)
	copier.CopyWithOption(&ret, &expr, copier.Option{IgnoreEmpty: false, DeepCopy: true})

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
			return ret
		}
	}

	return ret
}

func genText(ctx context.Context, expr *ast.CallExpr) []byte {
	return []byte(prettyPrint(ctx, expr))
}

func report(ctx context.Context, call *ssa.Call) {
	pass := ctx.Value(passKey).(*analysis.Pass)
	var actual, want string
	var gotActual, gotWant bool
	for _, v := range GetOperands(call) {
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

type contextKey int

const (
	_ contextKey = iota
	passKey
	ssaKey
	inspectorKey
	posMapKey
)

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
