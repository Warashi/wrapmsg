package a

import (
	"fmt"
)

func init() {
	t := T()

	// call method
	if err := t.Err(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.Err: %w"`
	}
	if err := t.Err(); err != nil {
		_ = fmt.Errorf("t.Err: %w", err)
	}

	// call method with line break
	if err := t.
		Err(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.Err: %w"`
	}
	if err := t.
		Err(); err != nil {
		_ = fmt.Errorf("t.Err: %w", err)
	}
}

func init() {
	t := T()

	// multi method chain
	if err := t.U().Err(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.U\.Err: %w"`
	}
	if err := t.U().Err(); err != nil {
		_ = fmt.Errorf("t.U.Err: %w", err)
	}

	// multi method chain with line break
	if err := t.
		U().
		Err(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "t\.U\.Err: %w"`
	}
	if err := t.
		U().
		Err(); err != nil {
		_ = fmt.Errorf("t.U.Err: %w", err)
	}
}
