package wrapmsg

import (
	"context"

	"github.com/Warashi/ssautil"
	"golang.org/x/tools/go/ssa"
)

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

func (w *walker) walkRefs(ctx context.Context, v ssautil.Referrerer) ([]string, bool) {
	for _, v := range ssautil.Referrers(v) {
		if r, ok := w.walk(ctx, v); ok {
			return r, true
		}
	}
	return nil, false
}

func (w *walker) walkOperands(ctx context.Context, v ssautil.Operander) ([]string, bool) {
	for _, v := range ssautil.Operands(v) {
		if r, ok := w.walk(ctx, v); ok {
			return r, true
		}
	}
	return nil, false
}

func (w *walker) walk(ctx context.Context, v ssautil.Poser) ([]string, bool) {
	if w.contains(v) {
		return nil, false
	}
	w.push(v)
	defer w.pop()

	switch v := v.(type) {
	case *ssa.Slice:
		return w.walkOperands(ctx, v)
	case *ssa.Alloc:
		return w.walkRefs(ctx, v)
	case *ssa.IndexAddr:
		return w.walkRefs(ctx, v)
	case *ssa.Store:
		return w.walkOperands(ctx, v)
	case *ssa.ChangeInterface:
		return w.walkOperands(ctx, v)
	case *ssa.Extract:
		return w.walkOperands(ctx, v)
	case *ssa.UnOp:
		return w.walkRefs(ctx, v.X)
	case *ssa.Call:
		return formatCall(ctx, v)
	}
	return nil, false
}
