package transpiler

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	itemType      int
	selectionType int
	charType      int
	parameterType int
)

const (
	VariableDeclaration itemType = iota
	ArrDeclaration
	FunctionDeclaration
	VariableAssignment
	ArrAssignment
	SelectionIf
	SelectionElseIf
	SelectionElse
	LoopStatement
	LoopBreakStatement
	ReturnStatement
	ScopeClose
	MacroItem
	Empty
)

const (
	If selectionType = iota
	ElseIf
	Else
)

const (
	letter charType = iota
	number
	underscore
	other
)

const (
	VariableParameter parameterType = iota
	ArrayParameter
)

type Variable struct {
	// where variables cannot be collections
	identifier string
	dataType   primitiveType
	mut        bool
}

type Declaration struct {
	v Variable
	e Expression
}

// cannot parse array literals into function because of type inference
type FunctionCall struct {
	functionName string
	parameters   []Expression
	arrays       []Array
	order        []parameterType
}

type Operator struct {
	// necessary for the ExpressionItem interface
	operator string
}

type IfStatement struct {
	statements []SelectionStatement
}

type SelectionStatement struct {
	condition     Expression // boolean expression
	selectionType selectionType
}

type Function struct {
	identifier  string
	parameters  []Variable
	arrays      []Array
	paramsOrder []parameterType
	returnType  primitiveType
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

func parseCharType(char byte) charType {
	switch char {
	case 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90: // uppercase letter
		return letter
	case 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122: // lowercase
		return letter
	case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
		return number
	case 95:
		return underscore
	default:
		return other
	}
}

func parseIdentifier(id string, lineNum int) string {
	last := len(id) - 1
	if id[last] != ':' { // last character must be colon for type annotation
		panic(fmt.Sprintf("Line %d: Name '%s' is invalid because the last character must be a colon for type annotation, but here it is '%s'", lineNum+1, id, string(id[last])))
	}
	// returns string if valid name, otherwise panics
	if !(parseCharType(id[0]) == letter) { // doesn't begin with uppercase or lowercase letter
		panic(fmt.Sprintf("Line %d: Name '%s' is invalid because it does not begin with a letter", lineNum+1, id))
	}

	for i := 0; i < last; i++ { // last character can be syntactic character
		if !(parseCharType(id[i]) == letter || parseCharType(id[i]) == number || parseCharType(id[i]) == underscore) { // character other than letters, number or underscore
			panic(fmt.Sprintf("Line %d: Name '%s' is invalid because it contains invalid character '%s'", lineNum+1, id, string(id[i])))
		}
	}

	if _, ok := allKeywords()[id]; ok {
		panic(fmt.Sprintf("Line %d: Name '%s' is invalid because it is a keyword in Stella", lineNum+1, id))
	}
	// no exit conditions triggered, so name must be valid
	return id[:len(id)-1] // identifier without colon
}

func parseVariableDeclaration(line string, lineNum int, currentScope *Scope) Declaration {
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

	if _, v := (*currentScope).vars[id]; v {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, f := (*currentScope).functions[id]; f {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, a := (*currentScope).arrays[id]; a {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	}

	expectedType := readType(words[typeIndex], lineNum)

	if expectedType == IO {
		panic(fmt.Sprintf("Line %d: variables cannot have data type IO", lineNum+1))
	}

	if words[equalsIndex] != "=" {
		panic(fmt.Sprintf("Line %d: expected token '=' after type annotation", lineNum))
	}

	var equalsCharIndex int // expression is everything after equals
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			equalsCharIndex = i
			break
		}
		// should always find it because of exit condition above
	}

	expression := line[equalsCharIndex+1:]
	exprFound := parseExpression(expression, lineNum, currentScope)

	if exprFound.dataType != expectedType {
		panic(fmt.Sprintf("Line %d: expected type %s because of type annotation, found type %s", lineNum, expectedType.String(), exprFound.dataType.String()))
	}

	v := Variable{
		identifier: id,
		dataType:   expectedType,
		mut:        mut,
	}

	(*currentScope).vars[id] = v

	return Declaration{
		v: v,
		e: exprFound,
	}
}

func parseExpression(expression string, lineNum int, currentScope *Scope) Expression {
	var parsed []string
	var currentItem string
	var bracketCount int

	// TODO: account for string literal while parsing tokens
	for i := 0; i < len(expression); i++ {
		switch expression[i] { // split on operators and spaces
		case '"':
			var stringLiteral string
			for j := i; j < len(expression); j++ {
				stringLiteral += string(expression[j])
				if expression[j] == '"' && j != i {
					i = j
					break
				}
				if j == len(expression)-1 {
					panic(fmt.Sprintf("Line %d: unterminated string literal in expression", lineNum+1))
				}
			}
			parsed = append(parsed, stringLiteral)
		case '(', '{':
			if len(currentItem) != 0 { // parse it as function call, error will arise in parseFunctionCall() if there is one
				// if currentItem is not empty it must be a function name
				fnBracketCount := 0
				for j := i; j < len(expression); j++ {
					currentItem += string(expression[j])
					switch expression[j] {
					case '(':
						fnBracketCount++
					case ')':
						fnBracketCount--
					}
					if fnBracketCount == 0 {
						i = j
						break
					}
				}
				if fnBracketCount != 0 {
					panic(fmt.Sprintf("Line %d: bracket opened but never closed", lineNum+1))
				}
				parsed = append(parsed, currentItem)
				currentItem = ""
				continue
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
			if expression[i+1] == '=' {
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

		if binaryOperator && token != "-" {
			checkBinaryOperator(token, previous, next, lineNum, currentScope)
		}
		if unaryOperator {
			if token == "-" {
				_, prevBinary := binaryOperators()[previous]
				if i == 0 || prevBinary {
					checkUnaryOperator(token, previous, next, lineNum, currentScope)
				} else {
					checkBinaryOperator(token, previous, next, lineNum, currentScope)
				}
			} else {
				checkUnaryOperator(token, previous, next, lineNum, currentScope)
			}
		}
		if !unaryOperator && !binaryOperator {
			checkValue(token, previous, next, lineNum, currentScope)
		}
	}
	T := expressionType(parsed, lineNum, currentScope)
	return Expression{parsed, T}
}

func checkValue(value, previous, next string, lineNum int, currentScope *Scope) {
	// check valid pattern, checks for unexpected token error
	if value == "" { // possible that empty string gets passed from checkBinaryOperator() or checkUnaryOperator()
		panic(fmt.Sprintf("Line %d: Expected value before operator", lineNum+1))
	}

	identifier := false

	for i := 0; i < len(value); i++ {
		if value[i] == '(' {
			fnCall := parseFunctionCall(value, lineNum, currentScope)
			if _, ok := currentScope.functions[fnCall.functionName]; ok {
				identifier = true
			}
		} else if value[i] == '[' {
			_ = parseArrayIndexing(value, lineNum, currentScope)
			// check for valid array indexing
			return
		}
	}

	if _, ok := currentScope.vars[value]; ok {
		identifier = true
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

func checkUnaryOperator(operator, previous, next string, lineNum int, currentScope *Scope) {
	if _, ok := binaryOperators()[previous]; !ok {
		if previous != "" && previous != "(" && previous != "{" {
			panic(fmt.Sprintf("Line %d: invalid token %s before unary operator %s", lineNum+1, previous, operator))
		}
	}
	checkValue(next, operator, "", lineNum, currentScope) // doesn't matter in this case what next actually is
}

func checkBinaryOperator(operator, previous, next string, lineNum int, currentScope *Scope) {
	if previous == "(" || previous == ")" || previous == "{" || previous == "}" {
		return
	} else if next == "(" || next == ")" || next == "{" || next == "}" {
		return
	}

	checkValue(previous, "", operator, lineNum, currentScope) // again doesn't matter what comes before value
	if _, ok := unaryOperators()[next]; !ok {
		checkValue(next, operator, "", lineNum, currentScope) // as above
	}
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
	// doesn't need to check if statements are syntactically valid
	// just determines whether ot not they are statements
	words := strings.Fields(line)
	if len(words) == 0 {
		return true
	}
	switch words[0] {
	case "if", "loop", "let", "}":
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

	for i := 0; i < len(line); i++ {
		if line[i] == '!' {
			// macro
			return true
		}
	}
	return assignment
}

func parseMultiLineExpression(lines []string, lineNum int, currentScope *Scope) Expression {
	// TODO: multi-line expressions inside multi-line expressions (maybe)

	varsCopy := make(map[string]Variable) // used to later restore currentScope.vars to original
	// so that when variable declarations are actually parsed they don't throw an already declared error

	for k, v := range (*currentScope).vars {
		varsCopy[k] = v
	}

	bracketCount := 0
	exprCount := 0
	exprLine := -1
	var expr string
	for n := lineNum; n < len(lines); n++ {
		line := lines[n]
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				bracketCount++
			} else if line[i] == '}' {
				bracketCount--
			}
		}
		if bracketCount == 0 {
			break
		}
		if bracketCount != 1 {
			continue
		}
		if getItemType(lines[n], n, currentScope) == VariableDeclaration {
			_ = parseVariableDeclaration(lines[n], n, currentScope)
		}
		if exprCount >= 1 {
			panic(fmt.Sprintf("Line %d: found dead code after expression in multi-line expression", n+1))
		}
		if isStatement(line) {
			continue
		}
		expr = line
		exprLine = n
		exprCount++
	}
	if exprLine == -1 {
		(*currentScope).vars = varsCopy
		return Expression{
			items:    []string{},
			dataType: IO,
		}
	}
	to_return := parseExpression(expr, exprLine, currentScope)
	(*currentScope).vars = varsCopy
	return to_return
}

func parseParameters(params string, lineNum int) ([]Variable, []Array, []parameterType) {
	fields := strings.Split(params, ",")
	if fields[0] == "" {
		return []Variable{}, []Array{}, []parameterType{}
	}
	// will remove the commas
	var variables []Variable
	var arrays []Array
	var paramTypes []parameterType

	for _, param := range fields {
		words := strings.Fields(param)
		if len(words) != 2 {
			panic(fmt.Sprintf("Line %d: invalid element in list of parameters", lineNum+1))
		}

		if words[0][len(words[0])-1] != ':' {
			panic(fmt.Sprintf("Line %d: the last character of the parameter declaration %s is not a colon ':', which is required for a type annotation of the parameter", lineNum+1, words[0]))
		}
		ident := parseIdentifier(words[0], lineNum)

		var isArr bool
		for i := 0; i < len(words[1]); i++ {
			if words[1][i] == '[' {
				isArr = true
			}
		}

		if isArr {
			arrT := parseArrayType(words[1], lineNum)
			newArr := Array{
				identifier: ident,
				dataType:   arrT,
				mut:        false,
			}
			arrays = append(arrays, newArr)
			paramTypes = append(paramTypes, ArrayParameter)
		} else {
			T := readType(words[1], lineNum)
			if T == IO {
				panic(fmt.Sprintf("Line %d: function parameters cannot have type IO", lineNum+1))
			}
			newP := Variable{
				identifier: ident,
				dataType:   T,
				mut:        false,
			}

			variables = append(variables, newP)
			paramTypes = append(paramTypes, VariableParameter)
		}

	}
	return variables, arrays, paramTypes
}

func parseFunction(lines []string, lineNum int, currentScope *Scope) Function {
	var allLines string
	for _, l := range lines {
		allLines += l
		allLines += "\n"
	}

	line := lines[lineNum]
	// var returnType primitiveType
	words := strings.Fields(line)
	if words[0] != "function" {
		panic("parseFunction() somehow called without function keyword")
	}

	identEnd := 0

	var id []byte
	for i := 0; i < len(words[1]); i++ {
		if words[1][i] == '(' {
			break
		}
		id = append(id, words[1][i])
	}

	identifier := parseIdentifier(string(id)+":", lineNum)

	if _, v := (*currentScope).vars[identifier]; v {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, f := (*currentScope).functions[identifier]; f {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, a := (*currentScope).arrays[identifier]; a {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	}

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
	parameters, arrays, order := parseParameters(pStr, lineNum)

	for _, p := range parameters {
		(*currentScope).vars[p.identifier] = p
	}

	for _, arr := range arrays {
		(*currentScope).arrays[arr.identifier] = arr
	}

	afterIdent := line[identEnd+1:]
	afterWords := strings.Fields(afterIdent)

	if len(afterWords) < 3 {
		panic(fmt.Sprintf("Line %d: Expected return type annotation '->' and equals sign '=' after function indentifier", lineNum+1))
	}

	if afterWords[0] != "->" {
		panic(fmt.Sprintf("Line %d: expected return type annotation with '->'", lineNum+1))
	}

	returnType := readType(afterWords[1], lineNum)
	if returnType == IO {
		if identifier != "main" {
			panic(fmt.Sprintf("Line %d: only the main() function can have return type IO", lineNum+1))
		}
	}

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

	var expression Expression

	if strings.Trim(lines[lineNum][exprStart:], " ")[0] != '{' { // single-line expression
		expression = parseExpression(lines[lineNum][exprStart:], lineNum, currentScope)
	} else {
		expression = parseMultiLineExpression(lines, lineNum, currentScope)
	}
	if expression.dataType != returnType {
		panic(fmt.Sprintf("Line %d: expected return type %v but found return type %v", lineNum+1, returnType, expression.dataType))
	}

	f := Function{
		parameters:  parameters,
		arrays:      arrays,
		paramsOrder: order,
		returnType:  returnType,
		identifier:  identifier,
	}

	(*currentScope).functions[f.identifier] = f

	return f
}

func parseFunctionCall(functionCall string, lineNum int, currentScope *Scope) FunctionCall {
	var ident, params string

	bracketCount := 0

	id := true
	for i := 0; i < len(functionCall); i++ {
		switch functionCall[i] {
		case '(':
			id = false
			bracketCount++
		case ')':
			bracketCount--
		}
		if id {
			ident += string(functionCall[i])
		} else {
			params += string(functionCall[i])
		}
	}

	if bracketCount != 0 {
		panic(fmt.Sprintf("Line %d: brackets opened but never closed", lineNum+1))
	}

	fn, ok := currentScope.functions[ident]

	if !ok {
		panic(fmt.Sprintf("Line %d: function %s not in scope", lineNum+1, ident+"()"))
	}

	if len(params) > 2 { // remove brackets in case of function with more than zero parameters
		params = params[1 : len(params)-1]
	}

	bracketCount = 0
	var parameterExprs []string
	var currentParam string

	for i := 0; i < len(params); i++ {
		switch params[i] {
		case '(':
			bracketCount++
			currentParam += string(params[i])
		case ')':
			bracketCount--
			currentParam += string(params[i])
		case ',':
			if bracketCount == 0 {
				parameterExprs = append(parameterExprs, currentParam)
				currentParam = ""
			}
		case ' ':
			if len(currentParam) != 0 {
				currentParam += string(params[i])
			}
		default:
			currentParam += string(params[i])
		}

		if i == len(params)-1 {
			parameterExprs = append(parameterExprs, currentParam)
			break
		}
	}

	if len(parameterExprs) != len(fn.paramsOrder) {
		panic(fmt.Sprintf("Line %d: function %s takes %d arguments but %d were given", lineNum+1, fn.identifier, len(fn.parameters), len(parameterExprs)))
	}

	var parameterExpressions []Expression
	var arrays []Array

	var variableCount, arrayCount int

	for i := 0; i < len(fn.paramsOrder); i++ {
		if fn.paramsOrder[i] == VariableParameter {
			expression := parseExpression(parameterExprs[i], lineNum, currentScope)
			if expression.dataType != fn.parameters[variableCount].dataType {
				panic(fmt.Sprintf("Line %d: cannot use expression of type %v as argument of type %v", lineNum+1, expression.dataType.String(), fn.parameters[i].dataType.String()))
			}
			parameterExpressions = append(parameterExpressions, expression)
			variableCount++
		} else {
			arrayExpression := parseArrayExpression(parameterExprs[i], lineNum, currentScope)
			if arrayExpression.dataType.baseType == fn.arrays[arrayCount].dataType.baseType {
				if len(arrayExpression.dataType.dimensions) != len(fn.arrays[arrayCount].dataType.dimensions) {
					panic(fmt.Sprintf("Line %dL expression does not have same number of dimensions as array parameter", lineNum+1))
				}
				for i := 0; i < len(arrayExpression.dataType.dimensions); i++ {
					if arrayExpression.dataType.dimensions[i] != fn.arrays[arrayCount].dataType.dimensions[i] {
						panic(fmt.Sprintf("Line %d: expression does not have same dimension size as array parameter", lineNum+1))
					}
				}
			} else {
				panic(fmt.Sprintf("Line %d: expression does not have same base type as array parameter", lineNum+1))
			}
			arrays = append(arrays, arrayExpression)
			arrayCount++
		}
	}

	return FunctionCall{
		functionName: ident,
		parameters:   parameterExpressions,
		arrays:       arrays,
		order:        fn.paramsOrder,
	}
}

func parseIfStatement(lineNum int, lines []string, currentScope *Scope) IfStatement {
	first := parseSelection(lineNum, lines, currentScope)

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
				next := parseSelection(lineNum, lines, currentScope)

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

func parseSelection(lineNum int, lines []string, currentScope *Scope) SelectionStatement {
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

	exprStart := 0
	var currentWord string

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case ' ':
			currentWord = ""
		default:
			currentWord += string(line[i])
		}

		if currentWord == "if" || currentWord == "else" {
			exprStart = i + 1
		}
	}

	exprEnd := 0

	for i := len(line) - 1; i > 0; i-- {
		if line[i] == '{' {
			exprEnd = i
		}
	}

	var condition Expression
	var expr string

	if exprStart+1 < exprEnd {
		expr = line[exprStart:exprEnd]
	}

	if T == Else {
		if len(strings.Fields(expr)) != 0 {
			panic(fmt.Sprintf("Line %d: else statements cannot contain a condition", lineNum+1))
		}
		condition = Expression{
			items:    []string{},
			dataType: Bool,
		}
	} else {
		condition = parseExpression(expr, lineNum, currentScope)
	}

	if condition.dataType != Bool {
		panic(fmt.Sprintf("Line %d: if statement found with non-boolean condition", lineNum+1))
	}

	return SelectionStatement{
		selectionType: T,
		condition:     condition,
	}
}

func parseAssignment(lines []string, lineNum int, currentScope *Scope) Assignment {
	line := lines[lineNum]
	words := strings.Fields(line)

	if len(words) < 3 {
		panic(fmt.Sprintf("Line %d: invalid assignment", lineNum+1))
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
			exprStart = i + 1
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

func getScopeType(lines []string, lineNum int) ScopeType {
	line := lines[lineNum]
	words := strings.Fields(line)
	switch words[0] {
	case "function":
		return FunctionScope
	case "if", "else":
		return SelectionScope
	default:
		panic(fmt.Sprintf("Line %d: invalid opening of scope", lineNum+1))
	}
}

func declarationType(line string, lineNum int) itemType {
	words := strings.Fields(line)
	identifierIndex := 1
	if len(words) == 1 {
		panic(fmt.Sprintf("Line %d: array declaration on line with only let keyword", lineNum+1))
	}
	if words[1] == "mut" {
		identifierIndex = 2
	}
	typeIndex := identifierIndex + 1
	typeWord := words[typeIndex]

	for i := 0; i < len(typeWord); i++ {
		if typeWord[i] == '[' {
			return ArrDeclaration
		}
	}
	return VariableDeclaration
}

func assignmentType(line string, lineNum int, currentScope *Scope) itemType {
	words := strings.Fields(line)
	if len(words) == 0 {
		panic("assignmentType() called on empty line %d")
	}
	_, isVar := (*currentScope).vars[words[0]]
	_, isArray := (*currentScope).arrays[words[0]]
	if isVar {
		return VariableAssignment
	} else if isArray {
		return ArrAssignment
	}
	panic(fmt.Sprintf("Line %d: assignment to variable %s not in scope", lineNum+1, words[0]))
}

func getItemType(line string, lineNum int, currentScope *Scope) itemType {
	words := strings.Fields(line)
	if len(words) == 0 {
		return Empty
	}
	switch words[0] {
	case "function":
		return FunctionDeclaration
	case "let":
		return declarationType(line, lineNum)
	case "if":
		return SelectionIf
	case "loop":
		return LoopStatement
	case "break", "continue":
		return LoopBreakStatement
	case "}":
		if len(words) == 1 {
			return ScopeClose
		} //-> must be at least 2
		if words[1] != "else" {
			panic(fmt.Sprintf("Line %d: only an else/else if statement can be opened on the same line where another scope is closed", lineNum+1))
		}
		if len(words) < 3 {
			panic(fmt.Sprintf("Line %d: keyword else followed by nothing", lineNum+1))
		}

		if words[2] == "if" {
			return SelectionElseIf
		}
		return SelectionElse
	default:
		if len(words) == 1 && words[0] == "}" {
			return ScopeClose
		}
		if !isStatement(line) {
			return ReturnStatement
		}

		for i := 0; i < len(line); i++ {
			if line[i] == '!' {
				return MacroItem
			}
		}

		for _, word := range words {
			if word == "=" {
				return assignmentType(line, lineNum, currentScope)
			}
		}

	}
	panic(fmt.Sprintf("Line %d: invalid line", lineNum+1))
	// shouldn't even be possible to get this
}

func parseScopeCloser(lines []string, lineNum int) ScopeCloser { // greatest function of all time
	line := lines[lineNum]
	words := strings.Fields(line)
	if len(words) != 1 {
		panic("parseScopeCloser() somehow called with line length != 1 🐐")
	}

	if words[0] != "}" {
		panic("parseScopeCloser() somehow called with words[0] != } 🐐")
	}

	return ScopeCloser{
		closer: "}",
	}
}

func typeOfItem(item Transpileable) string {
	typeof := fmt.Sprintf("%v", reflect.TypeOf(item))
	afterDot := false

	var T string
	for i := 0; i < len(typeof); i++ {
		if afterDot {
			T += string(typeof[i])
		}
		if typeof[i] == '.' {
			afterDot = true
		}
	}
	return T
}

func parseScope(lines []string, lineNum int, scopeType ScopeType, parent *Scope) Scope {
	newScope := Scope{
		vars:      make(map[string]Variable),
		arrays:    make(map[string]Array),
		functions: make(map[string]Function),
		scopeType: scopeType,
		items:     []Transpileable{},
		parent:    parent,
	}

	if parent != nil {
		// manually copy as maps are reference types
		newScope.vars = make(map[string]Variable)
		for k, v := range (*parent).vars {
			newScope.vars[k] = v
		}

		newScope.functions = make(map[string]Function)
		for k, v := range (*parent).functions {
			newScope.functions[k] = v
		}

		newScope.arrays = make(map[string]Array)
		for k, v := range (*parent).arrays {
			newScope.arrays[k] = v
		}
	}

	// NOTE: should be called inluding opening line
	scopeEnd := findScopeEnd(lines, lineNum)

	start, end := lineNum+1, scopeEnd
	if scopeType == Global {
		start, end = lineNum, len(lines)
	}

	var bracketCount int
	// where target is the bracketCount required to be in main scope

	for n := start; n < end; n++ {

		line := lines[n]
		var inMainScope bool // whether or not it is inside the main scope being read
		for i := 0; i < len(line); i++ {
			switch line[i] {
			case '{':
				bracketCount++
			case '}':
				bracketCount--
			}

			if bracketCount == 0 {
				inMainScope = true
			}
		}

		T := getItemType(line, n, &newScope)
		if !inMainScope && T != ScopeClose {
			continue
		}

		switch T {
		case VariableDeclaration:
			declaration := parseVariableDeclaration(line, n, &newScope)
			if newScope.scopeType == Global {
				panic(fmt.Sprintf("Line %d: global variables are not allowed in Stella", n))
			}
			newScope.items = append(newScope.items, declaration)

		case ArrDeclaration:
			declaration := parseArrayDeclaration(line, n, &newScope)
			newScope.items = append(newScope.items, declaration)

		case FunctionDeclaration:
			subScope := Scope{}

			// copy manually as maps are reference types
			subScope.vars = make(map[string]Variable)
			for k, v := range newScope.vars {
				subScope.vars[k] = v
			}

			subScope.functions = make(map[string]Function)
			for k, v := range newScope.functions {
				subScope.functions[k] = v
			}

			subScope.arrays = make(map[string]Array)
			for k, v := range newScope.arrays {
				subScope.arrays[k] = v
			}

			fn := parseFunction(lines, n, &subScope)
			newScope.functions[fn.identifier] = fn
			newScope.items = append(newScope.items, fn)

			subScope = parseScope(lines, n, FunctionScope, &subScope)
			// kinda scuffed but I don't think this causes any problems
			newScope.items = append(newScope.items, subScope)
			ended := findScopeEnd(lines, n)
			n = ended - 1

		case VariableAssignment:
			assignment := parseAssignment(lines, n, &newScope)
			if newScope.scopeType == Global {
				panic(fmt.Sprintf("Line %d: global variables are not allowed in Stella", n))
			}
			newScope.items = append(newScope.items, assignment)

		case ArrAssignment:
			assignment := parseArrayAssignment(lines[n], n, &newScope)
			newScope.items = append(newScope.items, assignment)

		case ReturnStatement:
			if scopeType != FunctionScope {
				panic(fmt.Sprintf("Line %d: Found return statement outside function scope", n+1))
			}
			expr := parseExpression(line, n, &newScope)
			newScope.items = append(newScope.items, expr)

		case SelectionIf:

			subScope := Scope{}

			// copy manually as maps are reference types
			subScope.vars = make(map[string]Variable)
			for k, v := range newScope.vars {
				subScope.vars[k] = v
			}

			subScope.functions = make(map[string]Function)
			for k, v := range newScope.functions {
				subScope.functions[k] = v
			}

			subScope.arrays = make(map[string]Array)
			for k, v := range newScope.arrays {
				subScope.arrays[k] = v
			}

			ifStatement := parseSelection(n, lines, &subScope)
			newScope.items = append(newScope.items, ifStatement)

			subScope = parseScope(lines, n, SelectionScope, &newScope)
			newScope.items = append(newScope.items, subScope)
			ended := findScopeEnd(lines, n)
			n = ended - 1

		case SelectionElse, SelectionElseIf:
			if len(newScope.items) == 0 {
				panic(fmt.Sprintf("Line %d: else/else if statements must be preceded by other selection statements", n+1))
			}
			if typeOfItem(newScope.items[len(newScope.items)-1]) != "Scope" {
				panic(fmt.Sprintf("Line %d: else/else if statements must be preceded by other selection statements", n+1))
			}

			scopeCount := -1
			for i := n - 1; i >= 0; i-- {
				line := lines[i]
				for j := 0; j < len(line); j++ {
					switch line[j] {
					case '{':
						scopeCount++
					case '}':
						scopeCount--
					}
				}
				if scopeCount == 0 {
					if getItemType(lines[i], i, &newScope) != SelectionIf {
						panic(fmt.Sprintf("Line %d: else/else if statements must be preceded by if statements", n+1))
					}
					break
				}
			}

			subScope := Scope{}

			// copy manually as maps are reference types
			subScope.vars = make(map[string]Variable)
			for k, v := range newScope.vars {
				subScope.vars[k] = v
			}

			subScope.functions = make(map[string]Function)
			for k, v := range newScope.functions {
				subScope.functions[k] = v
			}

			subScope.arrays = make(map[string]Array)
			for k, v := range newScope.arrays {
				subScope.arrays[k] = v
			}

			ifStatement := parseSelection(n, lines, &subScope)
			newScope.items = append(newScope.items, ifStatement)

			subScope = parseScope(lines, n, SelectionScope, &newScope)
			newScope.items = append(newScope.items, subScope)
			ended := findScopeEnd(lines, n)
			n = ended - 1

		case LoopStatement:
			subScope := Scope{}

			// copy manually as maps are reference types
			subScope.vars = make(map[string]Variable)
			for k, v := range newScope.vars {
				subScope.vars[k] = v
			}

			subScope.functions = make(map[string]Function)
			for k, v := range newScope.functions {
				subScope.functions[k] = v
			}

			subScope.arrays = make(map[string]Array)
			for k, v := range newScope.arrays {
				subScope.arrays[k] = v
			}

			loop := parseLoop(lines[n], n, &subScope)
			newScope.items = append(newScope.items, loop)

			subScope = parseScope(lines, n, LoopScope, &newScope)
			newScope.items = append(newScope.items, subScope)
			ended := findScopeEnd(lines, n)
			n = ended - 1

		case LoopBreakStatement:
			b := parseBreak(lines[n], n)
			newScope.items = append(newScope.items, b)

			currentScope := newScope

			// check that break statement is inside at least one loop

			for {
				if currentScope.scopeType == Global {
					panic(fmt.Sprintf("Line %d: found break/continue statement not inside any loop", lineNum+1))
				} else if currentScope.scopeType == LoopScope {
					break
				}
				if currentScope.parent == nil {
					panic(fmt.Sprintf("Line %d: found break/continue statement not inside any loop", lineNum+1))
				} else if (*currentScope.parent).scopeType == LoopScope {
					break
				}
				currentScope = *currentScope.parent
			}

		case MacroItem:
			macro := parseMacro(lines[n], lineNum, &newScope)
			if newScope.scopeType == Global {
				panic(fmt.Sprintf("Line %d: found unexpected macro in global scope", n+1))
			}
			newScope.items = append(newScope.items, macro)

		case Empty:

		case ScopeClose:
			if n == end {
				closer := ScopeCloser{closer: "}"}
				newScope.items = append(newScope.items, closer)
			} else {
				closer := parseScopeCloser(lines, n)
				newScope.items = append(newScope.items, closer)

			}
		}
	}

	if len(newScope.items) == 0 {
		panic(fmt.Sprintf("Line %d: scope is empty", lineNum+1))
	}

	return newScope
}
