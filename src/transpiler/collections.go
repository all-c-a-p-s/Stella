package transpiler

import (
	"fmt"
	"strconv"
	"strings"
)

type Array struct {
	identifier string
	dataType   ArrayType
	mut        bool
}

// let nums: int[5] = [1, 2, 3, 4, 5]
type ArrayType struct {
	dimensions []int
	baseType   primitiveType
}

type BaseArray struct {
	values   []Expression
	dataType primitiveType
	length   int
}

type ArrayValue[T primitiveType] struct {
	children []*ArrayValue[T] // will be empty if it is a base-array
	elements []Expression     // will be empty if not a base-array
	baseType primitiveType
	length   int
}

type ArrayDeclaration struct {
	expr ArrayExpression
	arr  Array
}

type ArrayIndexing struct {
	arrayID  string
	dataType ArrayType
	index    Expression
}

type ArrayAssignment struct {
	expr ArrayExpression
	arr  Array
}

type ArrayIndexAssignment struct {
	arrIndex ArrayIndexing
	value    Expression
}

type ArrayExpression struct {
	stringValue string
	dataType    ArrayType
}

func squareBracketEnd(s string, start, lineNum int) int {
	// find index in line or square bracket close
	var bracketCount int
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '[':
			bracketCount++
		case ']':
			bracketCount--
		}
		if bracketCount == 0 {
			return i
		}
	}
	panic(fmt.Sprintf("line %d: square bracket opened but never closed", lineNum+1))
}

func parseArrayType(typeWord string, lineNum int) ArrayType {
	// parse array type based on type annotation
	squareBracketIndex := -1
	for i := 0; i < len(typeWord); i++ {
		if typeWord[i] == '[' {
			squareBracketIndex = i
			break
		}
	}

	if squareBracketIndex == -1 {
		panic("parseArrayType() called without square brackets in type")
	}

	baseT := typeWord[:squareBracketIndex]
	T := readType(baseT, lineNum)
	if T == IO {
		panic(fmt.Sprintf("Line %d: arrays cannot have the data type IO", lineNum+1))
	}

	dims := typeWord[squareBracketIndex:]
	bracketCount := 0

	// check for valid dimensions declaration
	var dimensions []int

	var currentNumStr string

	for i := 0; i < len(dims); i++ {
		switch dims[i] {
		case '[':
			if bracketCount != 0 {
				panic(fmt.Sprintf("Line %d: invalid square bracket opening in array type annotation", lineNum+1))
			}
			bracketCount++
		case ']':
			if bracketCount != 1 {
				panic(fmt.Sprintf("Line %d: invalid square bracket closing in array type annotation", lineNum+1))
			}
			n, err := strconv.Atoi(currentNumStr)
			if err != nil {
				panic(fmt.Sprintf("Line %d: failed to convert %s to integer in array type annotation", lineNum+1, currentNumStr))
			}

			dimensions = append(dimensions, n)
			currentNumStr = ""

			bracketCount--

		default:
			if _, ok := numbers()[string(dims[i])]; ok {
				currentNumStr += string(dims[i])
			} else {
				panic(fmt.Sprintf("Line %d: unexpected character %s in array type annotation", lineNum+1, string(dims[i])))
			}
		}
	}

	return ArrayType{
		baseType:   T,
		dimensions: dimensions,
	}
}

func parseBaseArray(arrayValue string, expectedType primitiveType, currentScope *Scope, lineNum int) BaseArray {
	// parses base array where the data type of the elements is a primitive type
	if len(arrayValue) < 2 {
		panic(fmt.Sprintf("line %d: length of array value cannot be less than two", lineNum+1))
	}
	if arrayValue[0] != '[' || arrayValue[len(arrayValue)-1] != ']' {
		panic("arrayValue passed into parseBaseArray() wasn't opened and closed with square brackets")
	}
	if len(arrayValue) == 2 {
		return BaseArray{
			values:   []Expression{},
			dataType: expectedType,
			length:   0,
		}
	}

	var currentElement string
	var elements []Expression

	var stringLiteralCount int

	for i := 1; i < len(arrayValue); i++ {
		if arrayValue[i] == '"' {
			if stringLiteralCount == 1 {
				stringLiteralCount = 0
			} else {
				stringLiteralCount = 1
			}
		}

		if stringLiteralCount != 0 {
			// inside string literal
			currentElement += string(arrayValue[i])
			continue
		}
		switch arrayValue[i] {
		case ',':
			expr := parseExpression(currentElement, lineNum, currentScope)
			if expr.dataType != expectedType {
				panic(fmt.Sprintf("Line %d: found element of type %v in array of type %v", lineNum+1, expr.dataType, expectedType))
			}
			elements = append(elements, expr)
			currentElement = ""
		case ' ':
			currentElement = ""
		case ']':
			expr := parseExpression(currentElement, lineNum, currentScope)
			if expr.dataType != expectedType {
				panic(fmt.Sprintf("Line %d: found element of type %v in array of type %v", lineNum+1, expr.dataType, expectedType))
			}
			elements = append(elements, expr)
			currentElement = ""
		default:
			currentElement += string(arrayValue[i])
		}
	}
	return BaseArray{
		values:   elements,
		dataType: expectedType,
		length:   len(elements),
	}
}

func parseArrayValue[T primitiveType](arrayValue string, expectedType primitiveType, currentScope *Scope, lineNum int) ArrayValue[T] {
	// parses value of either base or multi-dimensional array
	if len(arrayValue) < 2 {
		panic(fmt.Sprintf("line %d: length of array value cannot be less than two", lineNum+1))
	}
	if arrayValue[0] != '[' || arrayValue[len(arrayValue)-1] != ']' {
		panic("arrayValue passed into parseBaseArray() wasn't opened and closed with square brackets")
	}
	if len(arrayValue) == 2 {
		return ArrayValue[T]{
			children: []*ArrayValue[T]{},
			elements: []Expression{},
			baseType: expectedType,
			length:   0,
		}
	}

	var children []*ArrayValue[T]
	var elements []Expression
	var L int

	var bracketCount, maxBracketCount int
	// max bracket count of 1 means it is a base array

	for i := 0; i < len(arrayValue); i++ {
		switch arrayValue[i] {
		case '[':
			bracketCount++
		case ']':
			bracketCount--
		}
		maxBracketCount = max(maxBracketCount, bracketCount)
	}

	bracketCount = 0

	if maxBracketCount == 1 {
		// recursive base case
		base := parseBaseArray(arrayValue, expectedType, currentScope, lineNum)
		elements = base.values
		L = len(elements)
	} else {
		// multi-dimensional array
		for i := 0; i < len(arrayValue); i++ {
			switch arrayValue[i] {
			case '[':
				bracketCount++
			case ']':
				bracketCount--
			}
			if arrayValue[i] == '[' && bracketCount == 2 {
				end := squareBracketEnd(arrayValue, i, lineNum)
				var subArrayValue string
				if end == len(arrayValue)-1 {
					subArrayValue = arrayValue[i:]
				} else {
					subArrayValue = arrayValue[i : end+1]
				}
				child := parseArrayValue[T](subArrayValue, expectedType, currentScope, lineNum)
				// recusrive call on elements of current array
				children = append(children, &child)
			}
		}
		L = len(children)
	}
	return ArrayValue[T]{
		children: children,
		elements: elements,
		baseType: expectedType,
		length:   L,
	}
}

// let mut nums: int[5] = [1, 2, 3, 4, 5]
func parseArrayDeclaration(line string, lineNum int, currentScope *Scope) ArrayDeclaration {
	// TODO: multi-dimensional arrays

	var mut bool

	// check pattens declaration can match:
	words := strings.Fields(line)
	if len(words) == 0 {
		panic("parseArrayDeclaration() called on empty line")
	}
	if words[0] != "let" {
		panic(fmt.Sprintf("Line %d: array declaration without let keyword at beginning of line", lineNum+1))
	}
	identifierIndex := 1
	if len(words) == 1 {
		panic(fmt.Sprintf("Line %d: array declaration on line with only let keyword", lineNum+1))
	}
	if words[1] == "mut" {
		identifierIndex = 2
		mut = true
	}

	typeIndex, equalsIndex := identifierIndex+1, identifierIndex+2

	expectedType := parseArrayType(words[typeIndex], lineNum)

	id := parseIdentifier(words[identifierIndex], lineNum)

	if _, v := (*currentScope).vars[id]; v {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, f := (*currentScope).functions[id]; f {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	} else if _, a := (*currentScope).arrays[id]; a {
		panic(fmt.Sprintf("Line %d: %s already defined in this scope", lineNum+1, id))
	}

	if words[equalsIndex] != "=" {
		panic(fmt.Sprintf("Line %d: expected = sign but found %s", lineNum+1, words[equalsIndex]))
	}

	var equalsCharIndex int // expression is everything after equals
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			equalsCharIndex = i
			break
		}
		// should always find it because of exit condition above
	}

	if equalsCharIndex == len(line)-1 {
		panic(fmt.Sprintf("Line %d: found no value assigned to array %s in declaration statement", lineNum+1, id))
	}
	expression := strings.Trim(line[equalsCharIndex+1:], " ")
	arrFound := parseArrayExpression(expression, expectedType.baseType, lineNum, currentScope)

	if arrFound.dataType.baseType != expectedType.baseType {
		panic(fmt.Sprintf("Line %d: expected array of type %v found array of type %v", lineNum+1, expectedType.baseType, arrFound.dataType.baseType))
	}

	if arrFound.dataType.dimensions[0] != expectedType.dimensions[0] {
		panic(fmt.Sprintf("Line %d: expected array of length %d, found array of length %d", lineNum+1, arrFound.dataType.dimensions[0], expectedType.dimensions[0]))
	}

	arr := Array{
		mut:        mut,
		identifier: id,
		dataType:   expectedType,
	}

	(*currentScope).arrays[arr.identifier] = arr

	return ArrayDeclaration{
		arr:  arr,
		expr: arrFound,
	}
}

func parseArrayIndexing(indexing string, lineNum int, currentScope *Scope) ArrayIndexing {
	trimmed := strings.Trim(indexing, " ")
	if len(strings.Fields(trimmed)) > 1 {
		panic(fmt.Sprintf("Line %d: array indexing cannot contain a space", lineNum+1))
	}

	var squareBracketIndex int
Loop:
	for i := 0; i < len(trimmed); i++ {
		switch trimmed[i] {
		case '[':
			squareBracketIndex = i
			break Loop
		case ']':
			panic(fmt.Sprintf("Line %d: found closing bracket ] before opening bracket [ in array indexing", lineNum+1))
		}
	}

	id := trimmed[:squareBracketIndex]
	arr, ok := (*currentScope).arrays[id]

	if !ok {
		panic(fmt.Sprintf("Line %d: attempt to index array %s that is not in scope", lineNum+1, id))
	}

	dims := trimmed[squareBracketIndex:]
	bracketCount := 0

	var dimensions []Expression

	var currentNumStr string

	for i := 0; i < len(dims); i++ {
		switch dims[i] {
		case '[':
			if bracketCount != 0 {
				panic(fmt.Sprintf("Line %d: invalid square bracket opening in array indexing", lineNum+1))
			}
			bracketCount++
		case ']':
			if bracketCount != 1 {
				panic(fmt.Sprintf("Line %d: invalid square bracket closing in array indexing", lineNum+1))
			}
			expr := parseExpression(currentNumStr, lineNum, currentScope)

			// NOTE: this might need to be changed for multi-dimensional arrays
			if expr.dataType != Int {
				panic(fmt.Sprintf("Line %d: attempt to index arrays with expression evaluating to non-integer type %v", lineNum+1, expr.dataType))
			}

			dimensions = append(dimensions, expr)
			currentNumStr = ""

			bracketCount--

		default:
			currentNumStr += string(dims[i])
		}
	}

	if len(dimensions) > 1 {
		panic(fmt.Sprintf("Line %d: multi-dimensional array indexing is not currently supported", lineNum+1))
	} else if len(dimensions) == 0 {
		panic(fmt.Sprintf("Line %d: array indexing with no value", lineNum+1))
	}

	num, err := strconv.Atoi(dimensions[0].transpile())
	if err != nil { // integer literal -> we can check whether it is inside array bounds
		if num > arr.dataType.dimensions[0]-1 { // zero-indexed
			panic(fmt.Sprintf("Line %d: attempt to index element %d but array has size %d", lineNum+1, num, arr.dataType.dimensions[0]))
		}
	}

	return ArrayIndexing{
		arrayID:  id,
		dataType: arr.dataType,
		index:    dimensions[0],
	}
}

func parseArrayAssignment(line string, lineNum int, currentScope *Scope) ArrayAssignment {
	words := strings.Fields(line)
	if len(words) < 3 {
		panic(fmt.Sprintf("Line %d: invalid assignment", lineNum+1))
	}

	// patterns array assignment can match:
	arr, ok := (currentScope).arrays[words[0]]

	if !ok {
		panic(fmt.Sprintf("Line %d: first token of assignment does not match any arrays in scope", lineNum+1))
	}
	if !arr.mut {
		panic(fmt.Sprintf("Line %d: attempt to assign new value to immutable array %s", lineNum+1, arr.identifier))
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
		panic(fmt.Sprintf("Line %d: found no expression in assignment to variable %s", lineNum+1, arr.identifier))
	}

	expectedType := arr.dataType.baseType

	if exprStart == 0 || exprStart == len(line)-1 {
		panic(fmt.Sprintf("Line %d: found no expression in assignment to variable %s", lineNum+1, arr.identifier))
	}

	expr := line[exprStart:]

	arrayExpr := parseArrayExpression(expr, expectedType, lineNum, currentScope)

	if arrayExpr.dataType.baseType != expectedType {
		panic(fmt.Sprintf("Line %d: attempt tp assign value of base type %v to array of base type %v", lineNum+1, arrayExpr.dataType.baseType, expectedType))
	}

	if len(arrayExpr.dataType.dimensions) != len(arr.dataType.dimensions) {
		panic(fmt.Sprintf("Line %d: attempt to assign array value with %d dimensions to array with %d dimensions", lineNum+1, len(arrayExpr.dataType.dimensions), len(arr.dataType.dimensions)))
	}

	return ArrayAssignment{
		arr:  arr,
		expr: arrayExpr,
	}
}

func parseArrayIndexAssignment(line string, lineNum int, currentScope *Scope) ArrayIndexAssignment {
	// parse assignment to index of array
	words := strings.Fields(line)

	var identifier string

Loop:
	for i := 0; i < len(words[0]); i++ {
		switch words[0][i] {
		case '[':
			break Loop
		default:
			identifier += string(words[0][i])
		}
	}

	var leftSideType primitiveType
	arr, ok := (*currentScope).arrays[identifier]

	indexing := parseArrayIndexing(words[0], lineNum, currentScope)

	leftSideType = indexing.dataType.baseType
	if ok {
		if !arr.mut {
			panic(fmt.Sprintf("Line %d: attempt to assign new value to element of immutable array %s", lineNum+1, identifier))
		}
	} else {
		panic(fmt.Sprintf("Line %d: attempted assignment to array %s not in scope", lineNum+1, identifier))
	}

	var exprStart int

	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			exprStart = i + 1
			break
		}
	}

	if exprStart == 0 || exprStart == len(line)-1 {
		panic(fmt.Sprintf("Line %d: found no expression in assignment to variable %s", lineNum+1, identifier))
	}

	expr := line[exprStart:]
	rightSide := parseExpression(expr, lineNum, currentScope)
	if rightSide.dataType != leftSideType {
		// should panic in parsing anyway, but maybe I'll change something later and forget
		panic(fmt.Sprintf("Line %d: data type of right hand side of expression does not match data type of left hand side", lineNum+1))
	}

	return ArrayIndexAssignment{
		arrIndex: indexing,
		value:    rightSide,
	}
}

func parseArrayExpression(expr string, expectedType primitiveType, lineNum int, currentScope *Scope) ArrayExpression {
	// parses any array expression
	trimmed := strings.Trim(expr, " ")

	if trimmed[0] == '[' {
		// array literals
		T := parseBaseArray(expr, expectedType, currentScope, lineNum)
		return ArrayExpression{
			stringValue: expr,
			dataType: ArrayType{
				baseType:   T.dataType,
				dimensions: []int{T.length},
			},
		}
	}

	if arr, ok := (*currentScope).arrays[trimmed]; ok {
		return ArrayExpression{
			stringValue: expr,
			dataType:    arr.dataType,
		}
	}
	currentString := ""
Loop:
	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case '(':
			break Loop
		default:
			currentString += string(expr[i])
		}
	}
	if fn, ok := (*currentScope).functions[currentString]; ok {
		// functions returning derived type
		if fn.returnsDerived {
			return ArrayExpression{
				stringValue: expr,
				dataType:    fn.derivedReturnType,
			}
		}
	}
	panic(fmt.Sprintf("Line %d: invalid array expression", lineNum+1))
}

func parseMultiLineArrayExpression(lines []string, lineNum int, expectedType primitiveType, currentScope *Scope) ArrayExpression {
	// parses blocks evaluating to array expression
	varsCopy := make(map[string]Variable) // used to later restore currentScope.vars to original
	// so that when variable declarations are actually parsed they don't throw an already declared error
	arraysCopy := make(map[string]Array)

	for k, v := range (*currentScope).vars {
		varsCopy[k] = v
	}

	for k, v := range (*currentScope).arrays {
		arraysCopy[k] = v
	}
	// copy as reference types

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
		} else if getItemType(lines[n], n, currentScope) == ArrDeclaration {
			_ = parseArrayDeclaration(lines[n], n, currentScope)
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
		panic(fmt.Sprintf("Line %d: found no returned value if function", lineNum+1))
	}

	to_return := parseArrayExpression(expr, expectedType, exprLine, currentScope)
	(*currentScope).vars = varsCopy
	(*currentScope).arrays = arraysCopy

	return to_return
}

func findExpectedType(lines []string, lineNum int) primitiveType {
	// needs to loop backwards through lines to find function declartion with return type typeAnnotation
	// doesn't really need error checking as function declaration will already have been parsed
	for i := lineNum; i >= 0; i-- {
		words := strings.Fields(lines[i])
		if words[0] == "function" {
			typeAnnotation := words[len(words)-3] // not = or {
			expectedType := parseArrayType(typeAnnotation, i)
			return expectedType.baseType
		}
	}
	panic("in theory should never panic here lol")
}
