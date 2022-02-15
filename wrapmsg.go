package wrapmsg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/ssa"
)

const doc = "wrapmsg is ..."

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
	builtssa  *buildssa.SSA
	inspected *inspector.Inspector
	posMap    map[token.Pos]ast.Node
)

func printIndent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("\t")
	}
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

func getIdentName(v interface{ Pos() token.Pos }) []string {
	ident, ok := posMap[v.Pos()].(*ast.Ident)
	switch v := v.(type) {
	case *ssa.Slice:
		if ok {
			fmt.Println("Slice:", ident.Name)
		}
		return nil
	case *ssa.Alloc:
		if ok {
			fmt.Println("Alloc:", ident.Name)
		}
		return nil
	case *ssa.IndexAddr:
		if ok {
			fmt.Println("IndexAddr:", ident.Name)
		}
		return nil
	case *ssa.Store:
		if ok {
			fmt.Println("Store:", ident.Name)
		}
		return nil
	case *ssa.ChangeInterface:
		if ok {
			fmt.Println("ChangeInterface:", ident.Name)
		}
		return nil
	case *ssa.Call:
		if ok {
			fmt.Println("Call:", ident.Name)
		}
		break
	case *ssa.UnOp:
		if ok {
			fmt.Println("UnOp:", ident.Name)
		}
		break
	case *ssa.Parameter:
		fmt.Println("Parameter:", v.Object().Name())
		return []string{v.Object().Name()}
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
func (w *walker) walk(depth int, v poser) ([]string, bool) {
	printIndent(depth)
	fmt.Printf("%[1]v\t%[1]T\n", v)

	if w.contains(v) {
		return nil, false
	}
	w.push(v)
	defer w.pop()

	org := v
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
		return w.walkOperands(depth+1, v)
	case *ssa.UnOp:
		if r, ok := w.walk(depth+1, v.X); ok {
			return append(r, getIdentName(org)...), true
		}
		return getIdentName(v), true
	case *ssa.Parameter:
		return getIdentName(v), true
	case *ssa.Function:
		return getIdentName(v), true
	}
	return nil, false
}

func buildPosMap() {
	posMap = make(map[token.Pos]ast.Node)
	inspected.Preorder(nil, func(node ast.Node) {
		for i := node.Pos(); i <= node.End(); i++ {
			posMap[i] = node
		}
	})
}

func isErrorf(call *ssa.Call) bool {
	return call.Common().StaticCallee().Name() == "Errorf"
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

func getCallExpr(call *ssa.Call) *ast.CallExpr {
	for node := posMap[call.Pos()]; node != nil; node = posMap[node.End()+1] {
		if node, ok := node.(*ast.CallExpr); ok {
			return node
		}
	}
	return nil
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
	buf := new(bytes.Buffer)
	_ = format.Node(buf, token.NewFileSet(), expr)
	return buf.Bytes()
}

func report(pass *analysis.Pass, call *ssa.Call) {
	var actual, want string
	var gotActual, gotWant bool
	for _, v := range GetOperands(call) {
		fmt.Printf("%[1]v\t%[1]T\n", v)
		switch v := v.(type) {
		case *ssa.Const:
			if !gotActual {
				actual = constant.StringVal(v.Value)
				gotActual = true
			}
		case *ssa.Slice:
			w := new(walker)
			if r, ok := w.walk(0, v); ok {
				want = strings.Join(r, ".") + ": %w"
				gotWant = true
			}
		}
	}
	if gotWant && actual != want {
		node := getCallExpr(call)
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

func run(pass *analysis.Pass) (interface{}, error) {
	builtssa = pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	inspected = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	buildPosMap()
	for _, call := range iterateErrorf() {
		report(pass, call)
	}
	return nil, nil
}
