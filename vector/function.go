package vector

type function string

var functions = []function{
	"sin",
	"cos",
	"tan",
	"log",
	"ln",
}

func isFunc(str string) bool {
	for _, fn := range functions {
		if string(fn) == str {
			return true
		}
	}
	return false
}
