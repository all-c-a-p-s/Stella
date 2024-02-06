package transpiler

import (
	"fmt"
	"strings"
)

type Loop struct {
	condition Expression
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
