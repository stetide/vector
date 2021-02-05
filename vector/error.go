package vector

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

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
	return ""
}

// HelpErr is Help err
type HelpErr struct{}

func (e HelpErr) Error() string {
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
