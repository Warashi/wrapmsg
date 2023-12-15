package a

import (
	"errors"
	"fmt"
)

func myFunc() (err error) {
	defer func() {}()
	err = errors.New("my err")
	_ = fmt.Errorf("%w", err) // want `the error-wrapping message should be "errors\.New: %w"`
	return err
}
