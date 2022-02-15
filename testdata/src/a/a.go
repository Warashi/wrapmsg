package a

import (
	"context"
	"fmt"

	"a/b"
)

func f() error {
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
