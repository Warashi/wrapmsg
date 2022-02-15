package a

import (
	"fmt"
)

func f() error {
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
