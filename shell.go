package main

import (
	"andre/vector/vector"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func input(msg string) string {
	fmt.Print(msg)
	r := bufio.NewReader(os.Stdin)
	txt, _ := r.ReadString('\n')
	return txt
}

func push(a interface{}) {
	fmt.Println(">>", a)
}

func main() {
	if len(os.Args[1:]) > 0 {
		txt := strings.TrimSpace(strings.Join(os.Args[1:], " "))
		res, err := vector.Run(txt)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
		return
	}

	fmt.Println("VECTOR " + vector.VERSION)
	for {
		txt := input("$ ")
		txt = strings.TrimSpace(txt)
		if txt == "" {
			continue
		}

		lexer := vector.NewLexer(txt)
		tokens, err := lexer.GenerateTokens()
		if err != nil {
			push(err)
			continue
		}
		//fmt.Println(tokens)

		parser := vector.NewParser(tokens)
		ast, err := parser.Parse()
		if err != nil {
			switch err.(type) {
			case vector.ExitErr:
				// fmt.Println("VECTOR " + vector.VERSION)
				fmt.Println()
				return
			}
			push(err)
			continue
		}
		fmt.Println(ast)

		res, err := vector.Execute(ast)
		if err != nil {
			push(err)
			continue
		} else if res == nil {
			continue
		}
		push(res)
	}
}
