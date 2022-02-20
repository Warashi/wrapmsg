package a

import (
	"fmt"

	c "a/b"
)

func init() {
	// call other package with import alias
	if err := c.F(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.F: %w"`
	}
	if err := c.F(); err != nil {
		_ = fmt.Errorf("c.F: %w", err)
	}

	// call other package with import alias, line break
	if err := c.
		F(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.F: %w"`
	}
	if err := c.
		F(); err != nil {
		_ = fmt.Errorf("c.F: %w", err)
	}
}

func init() {
	// method chain other package with import alias
	if err := c.T().Err(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.T\.Err: %w"`
	}
	if err := c.T().Err(); err != nil {
		_ = fmt.Errorf("c.T.Err: %w", err)
	}

	// method chain other package with import alias, line break
	if err := c.
		T().
		Err(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.T\.Err: %w"`
	}
	if err := c.
		T().
		Err(); err != nil {
		_ = fmt.Errorf("c.T.Err: %w", err)
	}
}

func init() {
	// multi method chain other package with import alias
	if err := c.T().U().Err(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.T\.U\.Err: %w"`
	}
	if err := c.T().U().Err(); err != nil {
		_ = fmt.Errorf("c.T.U.Err: %w", err)
	}

	// multi method chain other package with import alias, line break
	if err := c.
		T().
		U().
		Err(); err != nil {
		_ = fmt.Errorf("hoge: %w", err) // want `the error-wrapping message should be "c\.T\.U\.Err: %w"`
	}
	if err := c.
		T().
		U().
		Err(); err != nil {
		_ = fmt.Errorf("c.T.U.Err: %w", err)
	}
}
