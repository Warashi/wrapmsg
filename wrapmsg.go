package wrapmsg

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
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

func genWrapmsg(posMap map[token.Pos]ast.Node, currentPackagePath string, call *ssa.Call) string {
	var prefix string
	ops := getOperands(call)
	for _, op := range ops {
		switch op := op.(type) {
		case *ssa.Call:
			prefix = genWrapmsg(posMap, currentPackagePath, op)
		case *ssa.UnOp:
			if ident, ok := posMap[op.X.Pos()].(*ast.Ident); ok {
				prefix = ident.Name
			}
		default:
		}
	}

	name := getCallName(call)
	if prefix != "" {
		// まだ再帰の途中
		return prefix + "." + name
	}
	// 再帰終わって最後のreturn
	pkg := getCallPackage(call)
	if currentPackagePath != pkg.Path() {
		// 現在のpackageと違うpackageを呼んでる
		return pkg.Name() + "." + name
	}
	return name
}

func getCallPackage(call *ssa.Call) *types.Package {
	return call.Common().StaticCallee().Package().Pkg
}

func getCallName(call *ssa.Call) string {
	return call.Common().StaticCallee().Name()
}

func getErrorf(instr ssa.Instruction) (*ssa.Call, bool) {
	call, ok := instr.(*ssa.Call)
	if !ok {
		return nil, false
	}

	if getCallName(call) != "Errorf" {
		return nil, false
	}

	return call, true
}

func getConstString(val ssa.Value) (string, bool) {
	msg, ok := val.(*ssa.Const)
	if !ok {
		return "", false
	}
	if msg.Value.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(msg.Value), true
}

func getOperands(v ssa.Instruction) []ssa.Value {
	ops := v.Operands(nil)
	r := make([]ssa.Value, 0, len(ops))
	for _, op := range ops {
		if *op == nil {
			break
		}
		r = append(r, *op)
	}
	return r
}

func GetWrapmsg(posMap map[token.Pos]ast.Node, pkg *types.Package, call *ssa.Call) (string, bool) {
	args := call.Common().Args
	values, ok := args[len(args)-1].(*ssa.Slice)
	if !ok {
		return "", false
	}
	alloc, ok := values.X.(*ssa.Alloc)
	if !ok {
		return "", false
	}

	for _, ref := range *alloc.Referrers() {
		if refs, ok := ref.(*ssa.IndexAddr); ok {
			for _, ref := range *refs.Referrers() {
				ops := getOperands(ref)
				for _, op := range ops {
					ci, ok := op.(*ssa.ChangeInterface)
					if !ok {
						continue
					}
					call, ok := ci.X.(*ssa.Call)
					if !ok {
						continue
					}
					funcNames := genWrapmsg(posMap, pkg.Path(), call)
					return funcNames + ": %w", true
				}
			}
		}
	}
	return "", false
}

func buildPosMap(inspect *inspector.Inspector) map[token.Pos]ast.Node {
	m := make(map[token.Pos]ast.Node)
	inspect.Preorder(nil, func(node ast.Node) {
		for i := node.Pos(); i <= node.End(); i++ {
			m[i] = node
		}
	})
	return m
}

func run(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	posMap := buildPosMap(inspect)

	for _, f := range s.SrcFuncs {
		fmt.Println(f)
		for _, b := range f.Blocks {
			fmt.Printf("\tBlock %d\n", b.Index)
			for _, instr := range b.Instrs {
				if call, ok := getErrorf(instr); ok {
					args := call.Common().Args
					wrapmsg, ok := getConstString(args[len(args)-2])
					if !ok {
						continue
					}
					want, ok := GetWrapmsg(posMap, s.Pkg.Pkg, call)
					if !ok {
						continue
					}
					if wrapmsg != want {
						pass.Reportf(call.Pos(), "wrapping error message should be %q", want)
					}
					/*
						if !isErrorType(lastArg) {
							continue
						}
					*/
				}
			}
		}
	}
	return nil, nil
}
