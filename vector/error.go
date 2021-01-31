package vector

// ExitErr is exit err
type ExitErr struct{}

func (e ExitErr) Error() string {
	return "exit"
}
