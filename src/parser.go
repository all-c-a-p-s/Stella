package main

import (
	"fmt"
	"strings"
)

type itemType int

const (
	TypeInt itemType = iota // doesn't matter what it is as long as it has a type
	TypeFloat
	TypeBool
	TypeByte
	TypeString
	UnaryOperator
	BinaryOperator
)

type Variable struct {
	// where variables cannot be collections
	identifier string
	dataType   primitiveType
	mut        bool
}

type FunctionCall struct {
	functionName string
	parameters   []Variable
}

type Operator struct {
	// necessary for the ExpressionItem interface
	operator string
}

type SelectionStatement struct {
	previous      *selectionStatement
	condition     string // boolean expression
	selectionType selectionType
}

type ExpressionItem struct {
	stringValue string
	itemType    itemType
}

type Function struct {
	parameters map[string]variable // use map to include variable name. This should be useful to translate to readable Go
	returnType primitiveType
}

type Expression struct {
	items    []ExpressionItem
	dataType primitiveType
}

func parseIdentifier(id string, lineNum int) string {
	// returns string if valid name, otherwise panics
	if !(parseCharType(id[0]) == letter) { // doesn't begin with uppercase or lowercase letter
		panic(fmt.Sprintf("Line %d: Name '%s' is invalid because it does not begin with a letter", lineNum+1, id))
	}

	last := len(id) - 1

	for i := 0; i < last; i++ { // last character can be syntactic character
		if !(parseCharType(id[i]) == letter || parseCharType(id[i]) == number || parseCharType(id[i]) == underscore) { // character other than letters, number or underscore
			panic(fmt.Sprintf("Line %d: Name '%s' is invalid because it contains invalid character '%s'", lineNum+1, id, string(id[i])))
		}
	}

	if id[last] != ':' { // last character must be colon for type annotation
		panic(fmt.Sprintf("Line %d: Name '%s' is invalid because the last character must be a colon for type annotation, but here it is '%s'", lineNum+1, id, string(id[last])))
	}

	if _, ok := allKeywords()[id]; ok {
		panic(fmt.Sprintf("Line %d: Name '%s' is invalid because it is a keyword in Stella", lineNum+1, id))
	}
	// no exit conditions triggered, so name must be valid
	return id[:len(id)-1] // identifier without colon
}

func parseVariableDeclaration(line string, lineNum int, currentScope *scope) Variable {
	// will be called after we are sure it is a variable that is being assigned
	var mut bool
	words := strings.Fields(line)
	if words[0] != "let" { // no idea how this function can evn be called without "let"
		panic(fmt.Sprintf("Line %d: Variable assignment without let keyword", lineNum+1))
	}
	identifierIndex := 1 // index where identifier is expected
	if words[1] == "mut" {
		identifierIndex = 2
		mut = true
	}

	typeIndex := identifierIndex + 1 // type must come after identifier
	equalsIndex := typeIndex + 1     //"=" must come after type

	id := parseIdentifier(words[identifierIndex], lineNum)
	typeIndex = int(readType(words[typeIndex], lineNum))

	expectedType := readType(words[typeIndex], lineNum)

	if words[equalsIndex] != "=" {
		panic(fmt.Sprintf("Line %d: expected token '=' after type annotation", lineNum))
	}

	var equalsCharIndex int // expression is everything after equals
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			equalsCharIndex = i
		}
		// should always find it because of exit condition above
	}

	expression := line[equalsCharIndex+1:]
	parsed := parseExpression(expression, lineNum, currentScope)
	exprType := expressionType(parsed, lineNum, currentScope)

	if exprType != expectedType {
		panic(fmt.Sprintf("Line %d: expected type %s because of type annotation, found type %s", lineNum, expectedType.String(), exprType.String()))
	}

	return Variable{
		identifier: id,
		dataType:   expectedType,
		mut:        mut,
	}
}

func parseExpression(expression string, lineNum int, currentScope *scope) []string {
	var parsed []string
	var currentItem string
	var bracketCount int
	// TODO: account for string literal while parsing tokens
	for i := 0; i < len(expression); i++ {
		switch expression[i] { // split on operators and spaces
		case '(', '{':
			// TODO: check if currentItem is the name of a function
			if currentItem != "" {
				parsed = append(parsed, currentItem)
			}
			bracketCount++
			currentItem = ""
			parsed = append(parsed, string(expression[i]))
		case ')', '}':
			if currentItem != "" {
				parsed = append(parsed, currentItem)
			}
			bracketCount--
			currentItem = ""
			parsed = append(parsed, string(expression[i]))
		case '+', '-', '*', '/': // can only be alone
			c := string(expression[i])
			if currentItem != "" {
				parsed = append(parsed, currentItem)
			}
			parsed = append(parsed, c)
			currentItem = ""
		case '&', '|', '=': // can only be 2 next to each other
			c := string(expression[i])
			if len(expression)-1 == i {
				panic(fmt.Sprintf("line %d: use of invalid operator %s in expression", lineNum+1, string(c)))
			}
			if string(expression[i+1]) == c { //
				if len(currentItem) != 0 {
					parsed = append(parsed, currentItem)
				}
				parsed = append(parsed, c+c)
				currentItem = ""
				i++ // skip next character because already added here
			} else {
				panic(fmt.Sprintf("Line %d: use of invalid operator '%s' in expression", lineNum+1, string(c)))
			}
		case '!', '<', '>': // can be alone or with another character
			c := string(expression[i])
			if i == len(expression)-1 {
				panic(fmt.Sprintf("Line %d: operator %s found at end of expression with no value after", lineNum+1, string(expression[i])))
			}
			if expression[i+1] == '=' { //
				if len(currentItem) != 0 {
					parsed = append(parsed, currentItem)
				}
				currentItem = ""
				parsed = append(parsed, c+"=")
				i++ // skip next character because already added here
			} else {
				if len(currentItem) != 0 {
					parsed = append(parsed, currentItem)
				}
				currentItem = ""
				parsed = append(parsed, c)
			}
		case ' ':
			if currentItem != "" {
				parsed = append(parsed, currentItem)
			}
			currentItem = ""
		default:
			currentItem += string(expression[i])
			if i == len(expression)-1 {
				currentItem += string(expression[i])
				parsed = append(parsed, currentItem)
			}
		}
	}
	if bracketCount != 0 {
		panic(fmt.Sprintf("Line %d: invalid brackets in expression", lineNum+1))
	}

	var previous, next string

	for i, token := range parsed {
		if i == 0 {
			previous = ""
		} else {
			previous = parsed[i-1]
		}

		if i == len(parsed)-1 {
			next = ""
		} else {
			next = parsed[i+1]
		}

		_, binaryOperator := binaryOperators()[token]
		_, unaryOperator := unaryOperators()[token]
		if token == "(" || token == ")" || token == "{" || token == "}" {
			continue
		}
		if binaryOperator {
			checkBinaryOperator(token, previous, next, lineNum, currentScope)
		} else if unaryOperator {
			checkUnaryOperator(token, previous, next, lineNum, currentScope)
		} else {
			checkValue(token, previous, next, lineNum, currentScope)
		}
	}
	return parsed
}

func checkValue(value, previous, next string, lineNum int, currentScope *scope) {
	// check valid pattern, checks for unexpected token error
	if value == "" { // possible that empty string gets passed from checkBinaryOperator() or checkUnaryOperator()
		panic(fmt.Sprintf("Line %d: Expected value before operator", lineNum+1))
	}
	identifier := true
	if _, ok := currentScope.vars[value]; !ok {
		if _, ok := currentScope.functions[value]; !ok {
			// TODO: fix this so that it actually fucking works lol
			identifier = false
		}
	}
	if !identifier {
		getValType(value, lineNum) // only used so this can panic in case of invalid token
	}
	_, binaryPrevious := binaryOperators()[previous]
	_, unaryPrevious := unaryOperators()[previous]

	if !binaryPrevious && !unaryPrevious {
		switch previous {
		case "(", "{", "":
		default:
			panic(fmt.Sprintf("Line %d: unexpected token %s before value %s", lineNum+1, previous, value))
		}
	}

	_, binaryNext := binaryOperators()[next]
	_, unaryNext := unaryOperators()[next]

	if !binaryNext && !unaryNext {
		switch next {
		case ")", "}", "":
		default:
			panic(fmt.Sprintf("Line %d: unexpected token %s after value %s", lineNum+1, next, value))
		}
	}
}

func checkUnaryOperator(operator, previous, next string, lineNum int, currentScope *scope) {
	if _, ok := binaryOperators()[previous]; !ok {
		if previous != "" && previous != "(" && previous != "{" {
			panic(fmt.Sprintf("Line %d: invalid token %s before unary operator %s", lineNum+1, previous, operator))
		}
	}
	checkValue(next, operator, "", lineNum, currentScope) // doesn't matter in this case what next actually is
}

// TODO: probably need a function to check brackets
func checkBinaryOperator(operator, previous, next string, lineNum int, currentScope *scope) {
	if previous == "(" || previous == ")" || previous == "{" || previous == "}" {
		return
	} else if next == "(" || next == ")" || next == "{" || next == "}" {
		return
	}

	checkValue(previous, "", operator, lineNum, currentScope) // again doesn't matter what comes before value
	checkValue(next, operator, "", lineNum, currentScope)     // as above
}

// patterns an expression can match:
// expression BINARY OPERATOR expression
// UNARY OPERATOR expression
// ( expression )
// { expression }

// patterns expression values can match
// function call
// variable
// UNARY OPERATOR variable
// integer/float literal
// string literal
// bool literal
// byte literal

// patters unary operator can match:
// BINARY OPERATOR UNARY OPERATOR value
// UNARY OPERATOR value

// if none of the above, panic("unexpected token")

// func parseFuctionDeclaration(lines []string, lineNum int)
