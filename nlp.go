package goincv

import (
	"go/parser"
	"go/token"
)

var fset = token.NewFileSet()

func IsGolangCodeFormatOK(src string) bool {
	_, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	return err == nil
}

func IsCodePunctuationInPairs(src string) bool {
	puns := []string{
		"[]",
		"{}",
		"()",
	}

	stack := &Stack{}

	nnum := 0

	for i := range src {
		if src[i] == '\n' {
			nnum++
		}
		for k := range puns {
			if src[i] == puns[k][1] {
				if stack.Len() > 0 && stack.Look() == puns[k][0] {
					stack.Pop()

					// log.Println(nnum, string(*stack))
					break
				} else {
					return false
				}
			}

			if src[i] == puns[k][0] {
				stack.Push(src[i])
				// log.Println(nnum, string(*stack))
				break
			}
		}

	}
	return stack.Len() == 0

}
