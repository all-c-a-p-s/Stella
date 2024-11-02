package transpiler

import (
	"fmt"
	"strconv"
	"strings"
)

type TupleExpressionType int

const (
	Literal TupleExpressionType = iota
	TupleVariable
	FnCall
)

// tuple implementaion in Go
/**
type tuple3[T0 any, T1 any, T2 any] struct {
	v0 T0
	v1 T1
	v2 T2
}

//need to generate tupleN depending on what size tuples the user uses

// type matching of tuple elements checked at compile time
*/

type Tuple struct {
	identifier string
	pattern    TuplePattern
	mut        bool
}

type TupleLiteral struct {
	values []Expression
}

type TuplePattern struct {
	dataTypes []primitiveType
}

type TupleDeclaration struct {
	t Tuple
	e TupleExpression
}

type TupleIndexing struct {
	t Tuple
	i int // integer because it needs to be checked at compile-time
}

type TupleExpression struct {
	fnCall   FunctionCall // optional
	literal  TupleLiteral // optional
	t        Tuple        // optional - one of three must be present
	exprType TupleExpressionType
}

type TupleAssignment struct {
	t Tuple
	e TupleExpression
}

// e.g. (string, int, int, float)
func parseTuplePattern(pattern string, lineNum int) TuplePattern {
	p := strings.Trim(pattern, " ")
	if len(p) < 2 {
		panic(fmt.Sprintf("Line %d: tuple pattern in invalid because it does not contain '()''", lineNum+1))
	}
	if !(p[0] == '(' && p[len(p)-1] == ')') {
		panic(fmt.Sprintf("Line %d: tuple pattern is invalid because it is not enclosed by parentheses", lineNum+1))
	}

	dataTypes := []primitiveType{}
	typeStrings := strings.Split(p[1:len(p)-1], ", ")
	for _, s := range typeStrings {
		T := readType(s, lineNum)
		dataTypes = append(dataTypes, T)
	}

	tupleImports = append(tupleImports, len(dataTypes))
	// later collected into hashset

	return TuplePattern{
		dataTypes: dataTypes,
	}
}

func matchTuplePattern(tuple TupleLiteral, pattern TuplePattern, lineNum int) struct{} {
	if len(tuple.values) != len(pattern.dataTypes) {
		panic(fmt.Sprintf("Line %d: tuple does not match expected tuple pattern because they do not have the same length", lineNum+1))
	}

	for i := 0; i < len(tuple.values); i++ {
		if tuple.values[i].dataType != pattern.dataTypes[i] {
			panic(fmt.Sprintf("Line %d: tuple does not match expected pattern because element %d has the wrong data type", lineNum+1, i+1))
		}
	}
	return struct{}{}
}

// e.g. ("seba", 16, 182, 66.5)
func parseTupleLiteral(tupleValue string, pattern TuplePattern, lineNum int, currentScope *Scope) TupleLiteral {
	trimmed := strings.Trim(tupleValue, " ")
	if !(trimmed[0] == '(' && trimmed[len(trimmed)-1] == ')') {
		panic(fmt.Sprintf("Line %d: tuple is invalid because it is not enclosed by brackets ()", lineNum+1))
	}

	elements := trimmed[1 : len(trimmed)-1]

	var elementStrings []string
	var stringLiteral bool
	var currentString string

	for i := 0; i < len(elements); i++ {
		if elements[i] == '"' {
			stringLiteral = !stringLiteral
		}
		switch elements[i] {
		case ',':
			if !stringLiteral {
				elementStrings = append(elementStrings, currentString)
				currentString = ""
			}
		default:
			currentString += string(elements[i])
		}
		if i == len(elements)-1 && len(currentString) != 0 {
			elementStrings = append(elementStrings, currentString)
		}
	}

	var expressions []Expression

	for _, s := range elementStrings {
		expr := parseExpression(s, lineNum, currentScope)
		expressions = append(expressions, expr)
	}

	tupleFound := TupleLiteral{
		values: expressions,
	}

	_ = matchTuplePattern(tupleFound, pattern, lineNum)

	return TupleLiteral{
		values: expressions,
	}
}

func parseTupleExpression(expr string, pattern TuplePattern, lineNum int, currentScope *Scope) TupleExpression {
	trimmed := strings.Trim(expr, " ")
	if len(trimmed) == 0 {
		panic(fmt.Sprintf("Line %d: tuple expression is empty", lineNum+1))
	}
	if trimmed[0] == '(' {
		literal := parseTupleLiteral(expr, pattern, lineNum, currentScope)
		// error handled by ^ if pattern doesn't match
		return TupleExpression{
			exprType: Literal,
			literal:  literal,
		}
	}
	var currentString string
Loop:
	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case '(':
			// know from above that it isn't first -> not a literal
			break Loop
		default:
			currentString += string(expr[i])
		}
	}
	currentString = strings.Trim(currentString, " ") // remove leading spaces
	if fn, ok := (*currentScope).functions[currentString]; ok {
		// functions returning derived type
		if fn.returnDomain == tuple {
			call := parseFunctionCall(strings.Trim(expr, " "), lineNum, currentScope)
			// check that function call is actually valid
			return TupleExpression{
				fnCall:   call,
				exprType: FnCall,
			}
		}
	}

	if len(strings.Fields(trimmed)) != 1 {
		panic(fmt.Sprintf("Line %d: invalid tuple expression", lineNum+1))
	}

	id := strings.Fields(trimmed)[0]
	t, ok := (*currentScope).tuples[id]
	if !ok {
		panic(fmt.Sprintf("Line %d: tuple %s not found in scope", lineNum+1, id))
	}

	return TupleExpression{
		t:        t,
		exprType: TupleVariable,
	}
}

func parseTupleDeclaration(line string, lineNum int, currentScope *Scope) TupleDeclaration {
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

	id := parseIdentifier(words[identifierIndex], lineNum)

	// TODO: add currentScope.tuples[] to all of these

	if _, v := (*currentScope).vars[id]; v {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, f := (*currentScope).functions[id]; f {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, a := (*currentScope).arrays[id]; a {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, t := (*currentScope).tuples[id]; t {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	}

	var colonIndex, equalsIndex int // character index

	var stringLiteral bool
	for i := 0; i < len(line); i++ {
		if line[i] == '"' {
			stringLiteral = !stringLiteral
		}

		if !stringLiteral {
			if line[i] == '=' {
				equalsIndex = i
				break
			} else if line[i] == ':' {
				colonIndex = i
			}
		}

		if i == len(line)-1 {
			panic(fmt.Sprintf("Line %d: found no equals sign in assignment to tuple", lineNum+1))
		}
	}

	expectedPattern := parseTuplePattern(line[colonIndex+1:equalsIndex], lineNum)

	expression := line[equalsIndex+1:]
	exprFound := parseTupleExpression(expression, expectedPattern, lineNum, currentScope)
	// ^ already checks that the types match

	t := Tuple{
		identifier: id,
		pattern:    expectedPattern,
		mut:        mut,
	}

	(*currentScope).tuples[id] = t

	return TupleDeclaration{
		t: t,
		e: exprFound,
	}
}

func parseTupleIndexing(indexing string, lineNum int, currentScope *Scope) TupleIndexing {
	var indexIndex int // cold variable name 🥶
	for i := 0; i < len(indexing); i++ {
		if indexing[i] == '.' {
			indexIndex = i
			break
		}

		if i == len(indexing)-1 {
			panic(fmt.Sprintf("Line %d: no index operator '.' found in tuple indexing", lineNum+1))
		}
	}

	id := indexing[:indexIndex]
	index := indexing[indexIndex+1:] // 🔥

	t, ok := (*currentScope).tuples[id]
	if !ok {
		panic(fmt.Sprintf("Line %d: tuple indexed %s not in current scope", lineNum+1, id))
	}

	i, err := strconv.Atoi(strings.Trim(index, " "))
	if err != nil {
		panic(fmt.Sprintf("Line %d: tuple index is invalid because it is not an integer literal", lineNum+1))
	}

	return TupleIndexing{
		t: t,
		i: i,
	}
}

func parseTupleAssignment(line string, lineNum int, currentScope *Scope) TupleAssignment {
	var equalsIndex int
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			equalsIndex = i
			break
		}
		if i == len(line)-1 {
			panic(fmt.Sprintf("Line %d: found no equals sign in tuple assignment", lineNum+1))
		}
	}

	id := strings.Trim(line[:equalsIndex], " ")
	t, ok := (*currentScope).tuples[id]

	if !t.mut {
		panic(fmt.Sprintf("Line %d: attempt to assign new value to immutable tuple %s", lineNum+1, id))
	}

	if !ok {
		panic(fmt.Sprintf("Line %d: assignment to tuple %s not in scope", lineNum+1, id))
	}

	if equalsIndex == len(line)-1 {
		panic(fmt.Sprintf("Line %d: no value assigned to %s in assignment", lineNum+1, id))
	}

	expr := parseTupleExpression(line[equalsIndex+1:], t.pattern, lineNum, currentScope)
	// ^ already checks that pattern matches

	return TupleAssignment{
		t: t,
		e: expr,
	}
}

func parseMultiLineTupleExpression(lines []string, lineNum int, pattern TuplePattern, currentScope *Scope) TupleExpression {
	// TODO: multi-line expressions inside multi-line expressions (maybe)

	varsCopy := make(map[string]Variable) // used to later restore currentScope.vars to original
	// so that when variable declarations are actually parsed they don't throw an already declared error
	arraysCopy := make(map[string]Array)
	tuplesCopy := make(map[string]Tuple)

	for k, v := range (*currentScope).vars {
		varsCopy[k] = v
	}
	for k, v := range (*currentScope).arrays {
		arraysCopy[k] = v
	}
	for k, v := range (*currentScope).tuples {
		tuplesCopy[k] = v
	}

	// manual copy as maps are reference types

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
			// ignore lines which are not in main scope
			continue
		}
		if getItemType(lines[n], n, currentScope) == VariableDeclaration {
			_ = parseVariableDeclaration(lines[n], n, currentScope)
		} else if getItemType(lines[n], n, currentScope) == ArrDeclaration {
			_ = parseArrayDeclaration(lines[n], n, currentScope)
		} else if getItemType(lines[n], n, currentScope) == TupDeclaration {
			_ = parseTupleDeclaration(lines[n], n, currentScope)
		}
		if exprCount >= 1 {
			if len(strings.Trim(line, " ")) > 0 {
				panic(fmt.Sprintf("Line %d: found dead code after expression in multi-line expression", n+1))
			} else {
				// blank lines are ok
				continue
			}
		}
		if isStatement(line) {
			continue
		}
		expr = line
		exprLine = n
		exprCount++
	}
	if exprLine == -1 {
		panic(fmt.Sprintf("Line %d: found no returned value in tuple block", lineNum+1))
	}
	to_return := parseTupleExpression(expr, pattern, exprLine, currentScope)
	(*currentScope).vars = varsCopy
	(*currentScope).arrays = arraysCopy
	(*currentScope).tuples = tuplesCopy
	// return maps to original
	return to_return
}

func isTupleIndexing(item string) bool { // helper function for Expression.transpile()
	if parseCharType(item[0]) == letter {
		var stringLiteral, byteLiteral bool
		for i := 0; i < len(item); i++ {
			// if the first syntactic character we find is '.', since we know the expression is valid
			// the expression must be tuple indexing
			if stringLiteral {
				if item[i] == 34 {
					stringLiteral = false
				}
				continue
			} else if byteLiteral {
				if item[i] == 39 {
					byteLiteral = false
				}
				continue
			}
			switch item[i] {
			case '.':
				return true
			case '(':
				return false
			case 34:
				stringLiteral = true
			case 39:
				byteLiteral = true
			}
		}
	}
	return false
}

func transpileTupleIndexing(item string) string {
	// we know the item is syntactically valid because it has already been checked
	var res string
	for i := 0; i < len(item); i++ {
		res += string(item[i])
		if item[i] == '.' {
			res += "v" + item[i+1:]
			break
		}
	}
	return res
}

func findExpectedPattern(lines []string, lineNum int) TuplePattern {
	// needs to loop backwards through lines to find function declartion with return type typeAnnotation
	// doesn't really need error checking as function declaration will already have been parsed
	var tuplePatternString string
	var opened bool
	for i := lineNum; i >= 0; i-- {
		words := strings.Fields(lines[i])
		if len(words) == 0 {
			continue
		}
		if words[0] == "function" {
			line := lines[i]
			var bracketCount int
		Loop:
			for i := len(line) - 1; i >= 0; i-- {
				if bracketCount == 0 && opened {
					// i.e. it has been opened and then closed
					break
				}
				switch line[i] {
				case ')':
					bracketCount--
					tuplePatternString = string(line[i]) + tuplePatternString
				case '(':
					bracketCount++
					tuplePatternString = string(line[i]) + tuplePatternString
					if bracketCount == 0 {
						break Loop
					}
				default:
					if opened {
						tuplePatternString = string(line[i]) + tuplePatternString
					}
				}
				if bracketCount != 0 {
					opened = true
				}

			}
		}
	}

	return parseTuplePattern(tuplePatternString, lineNum)
	// ^ shouldn't be possible for this to panic
}
