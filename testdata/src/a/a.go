package a

import (
	"context"
	"fmt"

	"a/b"
)

type ct struct{}

func (ct) Err(context.Context) error { return nil }

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

func (mm mmu) mmuerr(ctx context.Context) error {
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

func f() error {
	// multi method chain same package
	if err := T().U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.U\.Err: %w"`
	}
	if err := T().U().Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// method chain same package
	if err := T().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "T\.Err: %w"`
	}
	if err := T().Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// call method
	t := T()
	if err := t.Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.Err: %w"`
	}
	if err := t.Err(); err != nil {
		return fmt.Errorf("t.Err: %w", err)
	}

	// call other package
	if err := b.F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "b\.F: %w"`
	}
	if err := b.F(); err != nil {
		return fmt.Errorf("b.F: %w", err)
	}

	if err := ctx(context.Background()); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "ctx: %w"`
	}
	if err := ctx(context.Background()); err != nil {
		return fmt.Errorf("ctx: %w", err)
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

type ict interface {
	Err(context.Context) error
}
type mmu []ict
