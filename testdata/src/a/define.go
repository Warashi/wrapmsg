package a

import (
	"context"
)

type a struct {
	ct  ct
	ict ict
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
	Result() (struct{}, error)
}
type mmu []ict
