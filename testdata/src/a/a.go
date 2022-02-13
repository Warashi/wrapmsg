package a

import (
	"a/b"
	c "a/b"
	"context"
	"fmt"
)

func f() error {
	// call same package
	if err := g(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "g: %w"`
	}
	if err := g(); err != nil {
		return fmt.Errorf("g: %w", err)
	}

	if err := ctx(context.Background()); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "ctx: %w"`
	}
	if err := ctx(context.Background()); err != nil {
		return fmt.Errorf("ctx: %w", err)
	}

	tmp := context.Background()
	if err := ctx(tmp); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "ctx: %w"`
	}
	if err := ctx(tmp); err != nil {
		return fmt.Errorf("ctx: %w", err)
	}

	// method chain same package
	if err := T().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "T\.Err: %w"`
	}
	if err := T().Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// method chain same package with line break
	if err := T().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "T\.Err: %w"`
	}
	if err := T().
		Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// multi method chain same package
	if err := T().U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "T\.U\.Err: %w"`
	}
	if err := T().U().Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// multi method chain same package with line break
	if err := T().
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "T\.U\.Err: %w"`
	}
	if err := T().
		U().
		Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// method chain same package with args
	if err := T(1, 2).Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "T\.Err: %w"`
	}
	if err := T(1, 2).Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// method chain same package with args, line break
	if err := T(1, 2).
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "T\.Err: %w"`
	}
	if err := T(1, 2).
		Err(); err != nil {
		return fmt.Errorf("T.Err: %w", err)
	}

	// multi method chain same package with args
	if err := T(1, 2).U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "T\.U\.Err: %w"`
	}
	if err := T(1, 2).U().Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// multi method chain same package with args, line break
	if err := T(1, 2).
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "T\.U\.Err: %w"`
	}
	if err := T(1, 2).
		U().
		Err(); err != nil {
		return fmt.Errorf("T.U.Err: %w", err)
	}

	// call method
	t := T()
	if err := t.Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "t\.Err: %w"`
	}
	if err := t.Err(); err != nil {
		return fmt.Errorf("t.Err: %w", err)
	}

	// call method with line break
	if err := t.
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "t\.Err: %w"`
	}
	if err := t.
		Err(); err != nil {
		return fmt.Errorf("t.Err: %w", err)
	}

	// multi method chain
	if err := t.U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "t\.U\.Err: %w"`
	}
	if err := t.U().Err(); err != nil {
		return fmt.Errorf("t.U.Err: %w", err)
	}

	// multi method chain with line break
	if err := t.
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "t\.U\.Err: %w"`
	}
	if err := t.
		U().
		Err(); err != nil {
		return fmt.Errorf("t.U.Err: %w", err)
	}

	// call other package
	if err := b.F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "b\.F: %w"`
	}
	if err := b.F(); err != nil {
		return fmt.Errorf("b.F: %w", err)
	}

	// call other package with line break
	if err := b.
		F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "b\.F: %w"`
	}
	if err := b.
		F(); err != nil {
		return fmt.Errorf("b.F: %w", err)
	}

	// method chain other package
	if err := b.T().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "b\.T\.Err: %w"`
	}
	if err := b.T().Err(); err != nil {
		return fmt.Errorf("b.T.Err: %w", err)
	}

	// method chain other package with line break
	if err := b.
		T().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "b\.T\.Err: %w"`
	}
	if err := b.
		T().
		Err(); err != nil {
		return fmt.Errorf("b.T.Err: %w", err)
	}

	// multi method chain other package
	if err := b.T().U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "b\.T\.U\.Err: %w"`
	}
	if err := b.T().U().Err(); err != nil {
		return fmt.Errorf("b.T.U.Err: %w", err)
	}

	// multi method chain other package with line break
	if err := b.
		T().
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "b\.T\.U\.Err: %w"`
	}
	if err := b.
		T().
		U().
		Err(); err != nil {
		return fmt.Errorf("b.T.U.Err: %w", err)
	}

	// call other package with import alias
	if err := c.F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "c\.F: %w"`
	}
	if err := c.F(); err != nil {
		return fmt.Errorf("c.F: %w", err)
	}

	// call other package with import alias, line break
	if err := c.
		F(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "c\.F: %w"`
	}
	if err := c.
		F(); err != nil {
		return fmt.Errorf("c.F: %w", err)
	}

	// method chain other package with import alias
	if err := c.T().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "c\.T\.Err: %w"`
	}
	if err := c.T().Err(); err != nil {
		return fmt.Errorf("c.T.Err: %w", err)
	}

	// method chain other package with import alias, line break
	if err := c.
		T().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "c\.T\.Err: %w"`
	}
	if err := c.
		T().
		Err(); err != nil {
		return fmt.Errorf("c.T.Err: %w", err)
	}

	// multi method chain other package with import alias
	if err := c.T().U().Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "c\.T\.U\.Err: %w"`
	}
	if err := c.T().U().Err(); err != nil {
		return fmt.Errorf("c.T.U.Err: %w", err)
	}

	// multi method chain other package with import alias, line break
	if err := c.
		T().
		U().
		Err(); err != nil {
		return fmt.Errorf("hoge: %w", err) // want `wrapping error message should be "c\.T\.U\.Err: %w"`
	}
	if err := c.
		T().
		U().
		Err(); err != nil {
		return fmt.Errorf("c.T.U.Err: %w", err)
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
