package transpiler

import (
	"fmt"
	"strings"
)

type Loop struct {
	condition Expression
}

type BreakType int

const (
	Break BreakType = iota
	Continue
)

// break must be only token on the line
type BreakStatement struct {
	T BreakType
}

func parseLoop(line string, lineNum int, currentScope *Scope) Loop {
	words := strings.Fields(line)
	if words[0] != "loop" {
		panic("parseLoop() called without loop keyword")
	}
	trimmed := strings.Trim(line, " ")
	var exprEnd int

	if len(trimmed) == 4 {
		panic(fmt.Sprintf("Line %d: loop statement with blank line", lineNum+1))
	}

	for i := 0; i < len(trimmed); i++ {
		if trimmed[i] == '{' {
			if i != len(trimmed)-1 {
				panic(fmt.Sprintf("Line %d: scope opened in loop statement not at end of line", lineNum+1))
			}
			exprEnd = i
		}
	}

	expr := trimmed[4:exprEnd]
	expressionFound := parseExpression(expr, lineNum, currentScope)

	if expressionFound.dataType != Bool {
		panic(fmt.Sprintf("Line %d: use of loop statement without boolean condition", lineNum+1))
	}
	return Loop{
		condition: expressionFound,
	}
}

func parseBreak(line string, lineNum int) BreakStatement {
	words := strings.Fields(line)
	if len(words) != 1 {
		panic(fmt.Sprintf("Line %d: break statements must be the only token on the line", lineNum+1))
	}
	switch words[0] {
	case "break":
		return BreakStatement{
			T: Break,
		}
	case "continue":
		return BreakStatement{
			T: Continue,
		}
	default:
		panic(fmt.Sprintf("Line %d: found break statement with invalid keyword %s", lineNum+1, words[0]))
	}
}
