package a

import (
	"context"
	"fmt"

	"a/b"
	c "a/b"
)

func f() error {
	// call method of field
	tt := tt{}
	if err := tt.t.Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "tt.t.Err: %w"`
	}
	if err := tt.t.Err(); err != nil {
		return fmt.Errorf("tt.t.Err: %w", err)
	}
	if err := tt.ct.Err(context.Background()); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "tt.ct.Err: %w"`
	}
	if err := tt.ct.Err(context.Background()); err != nil {
		return fmt.Errorf("tt.ct.Err: %w", err)
	}

	// non-error
	_ = fmt.Errorf("new error")
	_ = fmt.Errorf("new error with format: %d", 10)
	var msg string
	_ = fmt.Errorf(msg)

	// call same package
	if err := g(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "g: %w"`
	}
	if err := g(); err != nil {
		return fmt.Errorf("g: %w", err)
	}

	if err := ctx(context.Background()); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "ctx: %w"`
	}
	if err := ctx(context.Background()); err != nil {
		return fmt.Errorf("ctx: %w", err)
	}

	tmp := context.Background()
	if err := ctx(tmp); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "ctx: %w"`
	}
	if err := ctx(tmp); err != nil {
		return fmt.Errorf("ctx: %w", err)
	}

	// range variable
	for _, ct := range []ct{{}, {}} {
		if err := ct.Err(context.Background()); err != nil {
			return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "ct.Err: %w"`
		}
	}
	for _, ct := range []ct{{}, {}} {
		if err := ct.Err(context.Background()); err != nil {
			return fmt.Errorf("ct.Err: %w", err)
		}
	}

	// method chain same package
	if err := T().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.Err: %w"`
	}
	if err := T().Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// method chain same package with line break
	if err := T().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.Err: %w"`
	}
	if err := T().
		Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// multi method chain same package
	if err := T().U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.U\.Err: %w"`
	}
	if err := T().U().Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// multi method chain same package with line break
	if err := T().
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.U\.Err: %w"`
	}
	if err := T().
		U().
		Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// method chain same package with args
	if err := T(1, 2).Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.Err: %w"`
	}
	if err := T(1, 2).Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// method chain same package with args, line break
	if err := T(1, 2).
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.Err: %w"`
	}
	if err := T(1, 2).
		Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// multi method chain same package with args
	if err := T(1, 2).U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.U\.Err: %w"`
	}
	if err := T(1, 2).U().Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// multi method chain same package with args, line break
	if err := T(1, 2).
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.U\.Err: %w"`
	}
	if err := T(1, 2).
		U().
		Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// call method
	t := T()
	if err := t.Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.Err: %w"`
	}
	if err := t.Err(); err != nil {
		return fmt.Errorf("t.Err: %w", err)
	}

	// call method with line break
	if err := t.
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.Err: %w"`
	}
	if err := t.
		Err(); err != nil {
		return fmt.Errorf("t.Err: %w", err)
	}

	// multi method chain
	if err := t.U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.U\.Err: %w"`
	}
	if err := t.U().Err(); err != nil {
		return fmt.Errorf("t.U.Err: %w", err)
	}

	// multi method chain with line break
	if err := t.
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.U\.Err: %w"`
	}
	if err := t.
		U().
		Err(); err != nil {
		return fmt.Errorf("t.U.Err: %w", err)
	}

	// call other package
	if err := b.F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "b\.F: %w"`
	}
	if err := b.F(); err != nil {
		return fmt.Errorf("b.F: %w", err)
	}

	// call other package with line break
	if err := b.
		F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "b\.F: %w"`
	}
	if err := b.
		F(); err != nil {
		return fmt.Errorf("b.F: %w", err)
	}

	// method chain other package
	if err := b.T().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "b\.T\.Err: %w"`
	}
	if err := b.T().Err(); err != nil {
		return fmt.Errorf("b.T.Err: %w", err)
	}

	// method chain other package with line break
	if err := b.
		T().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "b\.T\.Err: %w"`
	}
	if err := b.
		T().
		Err(); err != nil {
		return fmt.Errorf("b.T.Err: %w", err)
	}

	// multi method chain other package
	if err := b.T().U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "b\.T\.U\.Err: %w"`
	}
	if err := b.T().U().Err(); err != nil {
		return fmt.Errorf("b.T.U.Err: %w", err)
	}

	// multi method chain other package with line break
	if err := b.
		T().
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "b\.T\.U\.Err: %w"`
	}
	if err := b.
		T().
		U().
		Err(); err != nil {
		return fmt.Errorf("b.T.U.Err: %w", err)
	}

	// call other package with import alias
	if err := c.F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.F: %w"`
	}
	if err := c.F(); err != nil {
		return fmt.Errorf("c.F: %w", err)
	}

	// call other package with import alias, line break
	if err := c.
		F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.F: %w"`
	}
	if err := c.
		F(); err != nil {
		return fmt.Errorf("c.F: %w", err)
	}

	// method chain other package with import alias
	if err := c.T().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.T\.Err: %w"`
	}
	if err := c.T().Err(); err != nil {
		return fmt.Errorf("c.T.Err: %w", err)
	}

	// method chain other package with import alias, line break
	if err := c.
		T().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.T\.Err: %w"`
	}
	if err := c.
		T().
		Err(); err != nil {
		return fmt.Errorf("c.T.Err: %w", err)
	}

	// multi method chain other package with import alias
	if err := c.T().U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.T\.U\.Err: %w"`
	}
	if err := c.T().U().Err(); err != nil {
		return fmt.Errorf("c.T.U.Err: %w", err)
	}

	// multi method chain other package with import alias, line break
	if err := c.
		T().
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.T\.U\.Err: %w"`
	}
	if err := c.
		T().
		U().
		Err(); err != nil {
		return fmt.Errorf("c.T.U.Err: %w", err)
	}

	return nil
}

func (mm mmu) Err(ctx context.Context) error {
	for _, m := range mm {
		if err := m.Err(ctx); err != nil {
			return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "m.Err: %w"`
		}
		if err := m.Err(ctx); err != nil {
			return fmt.Errorf("m.Err: %w", err)
		}
	}
	return nil
}

type a struct {
	ct  ct
	ict interface{ Err(context.Context) error }
}

func (a *a) A(ctx context.Context) error {
	if err := a.ct.Err(ctx); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "a\.ct\.Err: %w"`
	}
	if err := a.ct.Err(ctx); err != nil {
		return fmt.Errorf("a.ct.Err: %w", err)
	}
	if err := a.ict.Err(ctx); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "a\.ict\.Err: %w"`
	}
	if err := a.ict.Err(ctx); err != nil {
		return fmt.Errorf("a.ict.Err: %w", err)
	}
	return nil
}

func g() error {
	return nil
}

func ctx(context.Context) error {
	return nil
}

func T(_ ...int) t {
	return t{}
}

type t struct{}

func (t) Err() error {
	return nil
}
func (t) U() t {
	return t{}
}

type tt struct {
	t  t
	ct ct
}
type itt struct {
	t  interface{ Err() error }
	ct interface{ Err(context.Context) error }
}

type ct struct{}

func (ct) Err(context.Context) error { return nil }

type r struct{}

func (r) r(context.Context) []*ct {
	return nil
}

type ict interface {
	Err(context.Context) error
}
type mmu []ict
