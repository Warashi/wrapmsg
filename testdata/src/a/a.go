package a

import (
	"context"
	"fmt"

	"a/b"
)

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
