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

type IfStatement struct {
	statements []SelectionStatement
}

type SelectionStatement struct {
	previous      *SelectionStatement
	condition     Expression // boolean expression
	selectionType selectionType
}

type Function struct {
	parameters []Variable
	returnType primitiveType
}

type Expression struct {
	items    []string
	dataType primitiveType
}

type Assignment struct {
	v Variable
	e Expression
}

type Iterator[T int64 | float64] struct {
	start    T
	step     T
	end      T
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
	exprFound := parseExpression(expression, lineNum, currentScope)

	if exprFound.dataType != expectedType {
		panic(fmt.Sprintf("Line %d: expected type %s because of type annotation, found type %s", lineNum, expectedType.String(), exprFound.dataType.String()))
	}

	return Variable{
		identifier: id,
		dataType:   expectedType,
		mut:        mut,
	}
}

func parseExpression(expression string, lineNum int, currentScope *scope) Expression {
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
	T := expressionType(parsed, lineNum, currentScope)
	return Expression{parsed, T}
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

func checkBrackets(bracket, previous, next string, lineNum int) {
	switch bracket {
	case "(", "{":
		switch previous {
		case ")", "}":
			panic(fmt.Sprintf("Line %d: invalid token %s found before bracket %s", lineNum+1, bracket, previous))
		}
		switch next {
		case ")", "}":
			panic(fmt.Sprintf("Line %d: invalid token %s found after bracket %s", lineNum+1, bracket, next))
		}
	case ")", "}":
		switch previous {
		case "(", "{":
			panic(fmt.Sprintf("Line %d: invalid token %s found before bracket %s", lineNum+1, bracket, previous))
		}
		switch next {
		case ")", "}":
			panic(fmt.Sprintf("Line %d: invalid token %s found after bracket %s", lineNum+1, bracket, next))
		}
	default:
		panic("checkBrackets() function somehow called without a bracket lmao")

	}
}

func isStatement(line string) bool {
	// doesn't need to check if statements are syntactically valis
	// just determines whether ot not they are statements
	words := strings.Fields(line)
	if len(words) == 0 {
		return true
	}
	switch words[0] {
	case "if", "else", "loop", "let":
		return true
	}
	assignment := false
	stringCount := 0
	for i := 0; i < len(line); i++ {
		if line[i] == '"' {
			if stringCount == 0 {
				stringCount++
			} else {
				stringCount--
			}
		}

		if line[i] == '=' && stringCount == 0 { // not inside string literal
			assignment = true
		}
	}
	return assignment
}

func parseMultiLineExpression(expression string, lineNum int, currentScope *scope) Expression {
	// TODO: multi-line expressions inside multi-line expressions (maybe)
	// NOTE: should be called not including brackets at beginning/end
	lines := strings.Split(expression, "\n")
	bracketCount := 0
	exprCount := 0
	var expr string
	for num, line := range lines {
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				bracketCount++
			} else if line[i] == '}' {
				bracketCount--
			}
		}
		if bracketCount != 0 {
			continue
		}
		if exprCount > 1 {
			panic(fmt.Sprintf("Line %d: found dead code after expression in multi-line expression", lineNum+num+1))
		}
		if isStatement(line) {
			continue
		}
		expr = line
		exprCount++
	}
	return parseExpression(expr, lineNum, currentScope)
}

func parseParameters(params string, lineNum int) []Variable {
	fields := strings.Split(params, ",")
	// will remove the commas
	var variables []Variable
	for _, param := range fields {
		words := strings.Fields(param)
		if len(words) != 2 {
			panic(fmt.Sprintf("Line %d: invalid element in list of parameters", lineNum+1))
		}

		if words[0][len(words[0])-1] != ':' {
			panic(fmt.Sprintf("Line %d: the last character of the parameter declaration %s is not a colon ':', which is required for a type annotation of the parameter", lineNum+1, words[0]))
		}
		ident := parseIdentifier(words[0], lineNum)
		T := readType(words[1], lineNum)
		newP := Variable{
			identifier: ident,
			dataType:   T,
			mut:        false,
		}
		variables = append(variables, newP)

	}
	return variables
}

func parseFunction(lines []string, lineNum int, currentScope *scope) Function {
	var allLines string
	for _, l := range lines {
		allLines += l
		allLines += "\n"
	}

	line := lines[lineNum]
	// var returnType primitiveType
	words := strings.Fields(line)
	if words[0] != "fn" {
		panic("parseFunction() somehow called without fn keyword")
	}

	identEnd := 0

	var paramsBytes []byte
	bracketCount := 0
	for i := 0; i < len(line); i++ {
		done := false
		switch line[i] {
		case '(':
			bracketCount++
		case ')':
			bracketCount--
			if bracketCount == 0 {
				done = true
				identEnd = i
			}
		}
		if bracketCount == 1 || done {
			paramsBytes = append(paramsBytes, line[i])
		}
		if done {
			break
		}
	}

	if identEnd == len(line) {
		panic(fmt.Sprintf("Line %d: expected return type annotation after function identifier", lineNum+1))
	}

	pStr := string(paramsBytes[1 : len(paramsBytes)-1])
	parameters := parseParameters(pStr, lineNum)

	afterIdent := line[identEnd+1:]
	afterWords := strings.Fields(afterIdent)

	if len(afterWords) < 3 {
		panic(fmt.Sprintf("Line %d: Expected return type annotation '->' and equals sign '=' after function indentifier", lineNum+1))
	}

	if afterWords[0] != "->" {
		panic(fmt.Sprintf("Line %d: expected return type annotation with '->'", lineNum+1))
	}

	returnType := readType(afterWords[1], lineNum)

	if afterWords[2] != "=" {
		panic(fmt.Sprintf("Line %d: expected equals sign '=' after return type annotation ->", lineNum+1))
	}

	exprStart := 0
	for i := 0; i < len(afterIdent); i++ {
		if afterIdent[i] == '=' {
			exprStart = identEnd + i + 2
		}
	}

	if exprStart == len(allLines) || exprStart == 0 {
		panic(fmt.Sprintf("Line %d: found no returned expression from function", lineNum+1))
	}

	expr := allLines[exprStart:]

	for i := 0; i < len(expr); i++ {
		if expr[i] == ' ' {
			continue
		} else {
			if expr[i] == '{' {
				bracketCount := 0
				bracketEnd := -1
				for j := i; j < len(expr); j++ {
					switch expr[j] {
					case '{':
						bracketCount++
					case '}':
						bracketCount--
					}
					if bracketCount == 0 {
						bracketEnd = j
					}
				}
				if bracketEnd == -1 || i == len(expr) {
					panic(fmt.Sprintf("Line %d: brackets in expression in function return never closed", lineNum+1))
				}
				expr = expr[i+1 : bracketEnd]
			}
			break
		}
	}

	expression := parseMultiLineExpression(expr, lineNum, currentScope)
	if expression.dataType != returnType {
		panic(fmt.Sprintf("Line %d: expected return type %v but found return type %v", lineNum+1, returnType, expression.dataType))
	}

	return Function{
		parameters: parameters,
		returnType: returnType,
	}
}

func parseIfStatement(lineNum int, lines []string, currentScope *scope) IfStatement {
	first := parseSelection(lineNum, lines, currentScope, nil)
	parent := &first

	statements := []SelectionStatement{first}

	for i := lineNum; i < len(lines); i++ {
		words := strings.Fields(lines[i])
		if len(words) == 0 {
			continue
		}
		if words[0] == "}" {
			if len(words) == 1 {
				break
			}
			if words[1] == "else" {
				next := parseSelection(lineNum, lines, currentScope, parent)
				parent = &next

				statements = append(statements, next)
			} else {
				panic(fmt.Sprintf("Line %d: expected either else or else if on same line as previous selection statement closed", lineNum+1))
			}
		}
	}

	return IfStatement{
		statements: statements,
	}
}

func parseSelection(lineNum int, lines []string, currentScope *scope, previous *SelectionStatement) SelectionStatement {
	line := lines[lineNum]
	words := strings.Fields(line)
	var T selectionType
	switch words[0] {
	case "if":
		T = If
	case "}": // opened on same line where previous selection statement closed
		if words[1] != "else" {
			panic(fmt.Sprintf("Line %d: expected else or else if after closed selection statement", lineNum+1))
		}
		if previous == nil {
			panic(fmt.Sprintf("Line %d: found else/else if statement without previous if statement", lineNum+1))
		}
		if len(words) == 2 {
			panic(fmt.Sprintf("Line %d: expected condition after keyword else", lineNum+1))
		}
		if words[2] == "if" {
			T = ElseIf
		} else {
			T = Else
		}
	default:
		panic("parseSelection() somehow called without if keyword or }")
	}

	if len(line) == 2 {
		panic(fmt.Sprintf("Line %d: if statement with no condition", lineNum+1))
	}

	exprEnd := 0

	for i := len(line) - 1; i > 0; i-- {
		if line[i] == '{' {
			exprEnd = i
		}
	}

	expr := line[2:exprEnd]
	condition := parseExpression(expr, lineNum, currentScope)

	if condition.dataType != Bool {
		panic(fmt.Sprintf("Line %d: if statement found with non-boolean condition", lineNum+1))
	}

	return SelectionStatement{
		previous:      previous,
		selectionType: T,
		condition:     condition,
	}
}

func parseAssignment(lines []string, lineNum int, currentScope *scope) Assignment {
	line := lines[lineNum]
	words := strings.Fields(line)

	if len(words) < 3 {
		panic(fmt.Sprintf("Line %d: invalid addignment", lineNum+1))
	}

	v, ok := (currentScope).vars[words[0]]

	if !ok {
		panic(fmt.Sprintf("Line %d: first token of assignment does not match any variables in current scope", lineNum+1))
	} else {
		if !v.mut {
			panic(fmt.Sprintf("Line %d: attempt to assign new value to immutable variable %s", lineNum+1, v.identifier))
		}
	}

	if words[1] != "=" {
		panic(fmt.Sprintf("Line %d: invalid assignment: equals sign must come directly after variable", lineNum+1))
	}

	var exprStart int

	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			exprStart = i
			break
		}
	}

	if exprStart == 0 || exprStart == len(line)-1 {
		panic(fmt.Sprintf("Line %d: found no expression in assignment to variable %s", lineNum+1, v.identifier))
	}

	expr := line[exprStart:]
	expression := parseExpression(expr, lineNum, currentScope)

	if expression.dataType != v.dataType {
		panic(fmt.Sprintf("Line %d: cannot assign epression of type %v to variable of type %v", lineNum+1, expression.dataType, v.dataType))
	}

	return Assignment{
		v: v,
		e: expression,
	}
}
