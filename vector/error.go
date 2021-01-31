package vector

// ExitErr is exit err
type ExitErr struct{}

func (e ExitErr) Error() string {
	return "exit"
}

// HelpErr is Help err
type HelpErr struct{}

func (e HelpErr) Error() string {
	return `HELP
Assign variable:   $ 'name' = 'expression'
End program:        $ quit | $ close | $ end | $ exit
Create vector:      vec('x' 'y' 'z' ...) | ['x' 'y' 'z' ...] | vec('x';'y';'z';...) | ['x';'y';'z';...]

Operator:
	Add:        '+'
	Subtract:   '-'
	Multiply:   '*'
	Devide:     '/' | ':'
	Power:      '^'
	Root:       '\'`
}
