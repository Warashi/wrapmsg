package a

import (
	"fmt"
	"testing"
)

func TestHoge(t *testing.T) {
	// non-error
	_ = fmt.Errorf("new error")
	_ = fmt.Errorf("new error with format: %d", 10)
	var msg string
	_ = fmt.Errorf(msg)
	t.Errorf("hoge")
}
