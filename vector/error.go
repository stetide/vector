package vector

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

// CharacterErr on invalid character
type CharacterErr struct{ msg string }

func (e CharacterErr) Error() string {
	return "Invalid character: " + e.msg
}

// SyntaxErr is syntax err
type SyntaxErr struct{ msg string }

func (e SyntaxErr) Error() string {
	// return "Invalid Syntax: " + e.msg
	return e.msg
}

// RuntimeErr from nodes
type RuntimeErr struct{ msg string }

func (e RuntimeErr) Error() string {
	return e.msg
}

// ImplementErr occurens when not implemented
type ImplementErr struct{ msg string }

func (e ImplementErr) Error() string {
	return e.msg
}

// ExitErr is exit err
type ExitErr struct{}

func (e ExitErr) Error() string {
	return "exit"
}

// ClearErr clears screen
type ClearErr struct{}

// Clear clears screen
func (e ClearErr) Clear() {
	switch runtime.GOOS {
	case "windows":
		log.Println("windows")
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func (e ClearErr) Error() string {
	switch runtime.GOOS {
	case "js":
		return "clear"
	default:
		return ""
	}
}

// HelpErr is Help err
type HelpErr struct{}

func (e HelpErr) Error() string {
	switch runtime.GOOS {
	case "js":
		return "help"
	default:
		return `HELP
Assign variable:    $ 'name' = 'expression'
End program:        $ quit | $ close | $ end | $ exit
Create vector:      $ vec('x' 'y' 'z' ...) | ['x' 'y' 'z' ...] | vec('x';'y';'z';...) | ['x';'y';'z';...]

Operator:
	Add:        '+'
	Subtract:   '-'
	Multiply:   '*'
	Devide:     '/' | ':'
	Power:      '^'
	Root:       '\'`
	}
}
