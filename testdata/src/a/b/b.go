package b

func F() error {
	return nil
}
func T(_ ...int) TT {
	return TT{}
}

type TT struct{}

func (TT) Err() error {
	return nil
}
func (TT) U() TT {
	return TT{}
}
