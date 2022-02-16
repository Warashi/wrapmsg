package wrapmsg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/printer"
	"go/token"
	"strconv"
	"strings"
	"sync"

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

var (
	pass      *analysis.Pass
	builtssa  *buildssa.SSA
	inspected *inspector.Inspector
	posMap    cMap
)

type cMap struct {
  base sync.Map
}

func (m *cMap) Load(key token.Pos) []ast.Node {
  r, _ := m.base.Load(key)
  if r == nil {
    return nil
  }
  return r.([]ast.Node)
}

func (m *cMap) Store(key token.Pos, value []ast.Node) {
  m.base.Store(key, value)
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

func getIdent(v poser) (*ast.Ident, bool) {
	for _, node := range posMap.Load(v.Pos()) {
		ident, ok := node.(*ast.Ident)
		if ok {
			return ident, true
		}
	}
	return nil, false
}

func getIdentName(v poser) []string {
	ident, ok := getIdent(v)
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

func (w *walker) walkRefs(depth int, v posReferrerer) ([]string, bool) {
	org := v
	for _, v := range *v.Referrers() {
		if r, ok := w.walk(depth, v); ok {
			return append(r, getIdentName(org)...), true
		}
	}
	return nil, false
}

func (w *walker) walkOperands(depth int, v posOperander) ([]string, bool) {
	org := v
	for _, v := range GetOperands(v) {
		if r, ok := w.walk(depth, v); ok {
			return append(r, getIdentName(org)...), true
		}
	}
	return nil, false
}

func format(expr ast.Expr) string {
	var b bytes.Buffer
	printer.Fprint(&b, pass.Fset, expr)
	return b.String()
}

func (w *walker) walk(depth int, v poser) ([]string, bool) {
	if w.contains(v) {
		return nil, false
	}
	w.push(v)
	defer w.pop()

	switch v := v.(type) {
	case *ssa.Const:
	case *ssa.Slice:
		return w.walkOperands(depth+1, v)
	case *ssa.Alloc:
		return w.walkRefs(depth+1, v)
	case *ssa.IndexAddr:
		return w.walkRefs(depth+1, v)
	case *ssa.Store:
		return w.walkOperands(depth+1, v)
	case *ssa.ChangeInterface:
		return w.walkOperands(depth+1, v)
	case *ssa.Call:
		call, ok := getCallExpr(v)
		if ok {
			return formatCallExpr(call), true
		}
	}
	return nil, false
}

func buildPosMap() {
	var mu sync.Mutex
	inspected.Preorder(nil, func(node ast.Node) {
		mu.Lock()
		defer mu.Unlock()
		for i := node.Pos(); i <= node.End(); i++ {
			stack := posMap.Load(i)
			stack = append(stack, node)
			posMap.Store(i, stack)
		}
	})
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

func iterateErrorf() []*ssa.Call {
	var r []*ssa.Call
	for _, f := range builtssa.SrcFuncs {
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

func getCallExpr(call poser) (*ast.CallExpr, bool) {
	for i := call.Pos(); i > 0; i-- {
		stack := posMap.Load(i)
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
	for i, arg := range expr.Args {
		c, ok := arg.(*ast.BasicLit)
		if !ok {
			continue
		}
		if c.Kind != token.STRING {
			continue
		}
		if strconv.Quote(actual) == c.Value && strconv.CanBackquote(actual) && c.Value == "`"+actual+"`" {
			continue
		}
		expr.Args[i] = &ast.BasicLit{
			ValuePos: c.ValuePos,
			Kind:     c.Kind,
			Value:    strconv.Quote(want),
		}
		return expr
	}
	return expr
}

func genText(expr *ast.CallExpr) []byte {
	return []byte(format(expr))
}

func report(pass *analysis.Pass, call *ssa.Call) {
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
			if r, ok := w.walk(0, v); ok && len(r) > 0 {
				want = strings.Join(r, ".") + ": %w"
				gotWant = true
			}
		}
	}
	if gotWant && actual != want {
		node, ok := getCallExpr(call)
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
				NewText: genText(replaceConst(node, actual, want)),
			}}}},
		})
	}
}

func run(p *analysis.Pass) (interface{}, error) {
	pass = p
	builtssa = pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	inspected = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	buildPosMap()
	for _, call := range iterateErrorf() {
		report(pass, call)
	}
	return nil, nil
}
