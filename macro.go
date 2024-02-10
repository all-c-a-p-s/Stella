package transpiler

import (
	"fmt"
	"strings"
)

type macroType int

const (
	Print macroType = iota
	Println
)

// println!("Hello world")

type Macro struct {
	value Expression
	T     macroType
}

func parseMacro(line string, lineNum int, currentScope *Scope) Macro {
	var macro string
	var bangIndex int
	for i := 0; i < len(line); i++ {
		if line[i] == '!' {
			bangIndex = i
			break
		}
		if i == len(line)-1 {
			panic("parsePrintStatement() called on line without ! macro")
		}
		macro += string(line[i])
	}
	macro = strings.Trim(macro, " ")
	var T macroType
	switch macro {
	case "print":
		T = Print
		imports = append(imports, "fmt")
	case "println":
		T = Println
		imports = append(imports, "fmt")
	default:
		panic(fmt.Sprintf("Line %d: attempt to use invalid macro %s!", lineNum+1, macro))
	}

	if bangIndex == len(line)-1 {
		panic(fmt.Sprintf("Line %d: attempt to call macro %s with no argument", lineNum+1, macro))
	}

	expr := parseExpression(line[bangIndex+1:], lineNum, currentScope)

	return Macro{
		T:     T,
		value: expr,
	}
}
