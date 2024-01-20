package main

// NOTE: I am currently completely rewriting this file in parser.go
// so it isn't really used

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	selectionType int
	charType      int
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

type variable struct { // name not used here because it will be stored as the key of the hashmap of variables
	dataType primitiveType
	mutable  bool
}

type function struct {
	parameters map[string]variable
	lineNum    int
	returnType primitiveType
}

type selectionStatement struct {
	previous      *selectionStatement // previous if statement e.g. if statement before else if
	condition     string
	selectionType selectionType
	begin         int // first line
	end           int // last line of statement, these are used to check if else if/else statements are opened on the same line as if statements are closed
}

type iterator struct {
	dataType primitiveType
	start    int
	end      int
	step     int
}

type forLoop struct {
	iterator iterator
	begin    int
	end      int
}

type infiniteLoop struct {
	exitCondition string
	begin         int
	end           int
}

type array struct {
	dataType   primitiveType // in case of multidimensional arrays, this will be the 'base' type
	vec        bool          // array or vector
	dimensions []int
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

func validName(name string, lineNum int) string {
	if !(parseCharType(name[0]) == letter) { // doesn't begin with uppercase or lowercase letter
		panic(fmt.Sprintf("Line %d: Name '%s' is invalid because it does not begin with a letter", lineNum+1, name))
	}

	for i := 0; i < len(name)-1; i++ { // last character can be syntactic character
		if !(parseCharType(name[i]) == letter || parseCharType(name[i]) == number || parseCharType(name[i]) == underscore) { // character other than letters, number or underscore
			panic(fmt.Sprintf("Line %d: Name '%s' is invalid because it contains invalid character '%s'", lineNum+1, name, string(name[i])))
		}
	}

	last := len(name) - 1

	if name[last] != ':' { // last character must be colon for type annotation
		panic(fmt.Sprintf("Line %d: Name '%s' is invalid because the last character must be a colon for type annotation, but here it is '%s'", lineNum+1, name, string(name[last])))
	}
	return name // no exit conditions triggered, so name must be valid
}

func syntacticCharacter(char byte) bool { // characters which have a syntactic function
	switch char {
	case ':', '(', ')', '[', ']', '{', '}':
		return true
	default:
		return false
	}
}

func comparisonOperator(operator string) bool {
	switch operator {
	case "==", "<=", ">=", "<", ">", "!=":
		return true
	default:
		return false
	}
}

func declarationKeyword(word string) bool {
	if word == "let" || word == "const" || word == "func" || word == "arr" || word == "vec" {
		return true
	}
	return false
}

func removeSyntacticChars(word string) (removed string) {
	for i := 0; i < len(word); i++ {
		if !syntacticCharacter(word[i]) {
			removed += string(word[i])
		}
	}
	return removed
}

func readName(lines []string, lineNum int) string {
	line := lines[lineNum]
	words := strings.Fields(line)

	for i := 0; i < len(words); i++ {
		if declarationKeyword(words[i]) && words[i] != "func" {
			if syntacticCharacter(words[i+1][len(words[i+1])-1]) { // names can have syntactic characters e.g. ':' or '(' after them without a space
				return validName(words[i+1][:len(words[i+1])], lineNum)
			}
			return validName(words[i+1], lineNum)
		} else if words[i] == "func" {
			var funcName string
			for j := 0; j < len(words[i+1]); j++ {
				if words[i+1][j] == '(' { // start of function parameters
					break
				}
				funcName += string(words[i+1][j])
			}
		}
	}
	panic(fmt.Sprintf("Line %d: invalid name in declaration", lineNum+1))
}

func readVariable(lines []string, lineNum int) (newVariable variable) {
	var dataType primitiveType
	mut := false

	words := strings.Fields(lines[lineNum])

	if words[1] == "mut" {
		mut = true
	}

	for i := 0; i < len(words); i++ {
		if words[i][len(words[i])-1] == ':' { // last character of colon indicates that this is the variable name and the type is next
			dataType = readType(words[i+1], lineNum)
		}
	}
	newVariable.dataType = dataType
	newVariable.mutable = mut
	return newVariable
}

/**
func readVariables(lines []string, scope *scope) {
	scopeCount := 0 // used to keep track of scopes opened/closed
	for lineNum, line := range lines[(*scope).begin:(*scope).end] {
		words := strings.Fields(line)
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++
			} else if line[i] == '}' {
				scopeCount--
			}
		}
		for _, word := range words {
			if word == "let" && scopeCount == 0 { // only read variables local to the current scope, not its subscopes
				(*scope).vars[readName(lines, lineNum)] = readVariable(lines, lineNum)
			}
		}
	}
}
*/

func readFunction(lines []string, lineNum int) (newFunction function) { // should only be called on a line once the func keyword has already been read
	var returnType primitiveType
	void := true

	words := strings.Fields(lines[lineNum])
	for i := 0; i < len(words); i++ {
		if words[i] == "->" {
			void = false
			returnType = readType(words[i+1], lineNum)
		}
	}

	if void {
		panic(fmt.Sprintf("Line %d: function has no return type", lineNum))
	}

	newFunction.returnType = returnType
	newFunction.parameters = readParameters(words[1], lineNum)

	return newFunction
}

func readFunctions(lines []string, scope *scope) {
	scopeCount := 0 // used to keep track of scopes opened/closed
	for lineNum, line := range lines[(*scope).begin:(*scope).end] {
		words := strings.Fields(line)
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++
			} else if line[i] == '}' {
				scopeCount--
			}
		}
		for _, word := range words {
			if word == "func" && scopeCount == 0 { // only read variables local to the current scope, not its subscopes
				(*scope).functions[readName(lines, lineNum)] = readFunction(lines, lineNum)
			}
		}
	}
}

func readParameters(funcName string, lineNum int) map[string]variable {
	parameters := map[string]variable{}
	bracketCount := 0
	params := ""
	for i := 0; i < len(funcName); i++ {
		if funcName[i] == '(' {
			bracketCount++
		} else if funcName[i] == ')' {
			bracketCount--
		}

		if bracketCount == 1 {
			params += string(funcName[i])
		}
	}

	params = removeSyntacticChars(params) // removes opening bracket
	funcParams := strings.Split(params, ",")

	for _, param := range funcParams {
		words := strings.Fields(param)
		if len(words) != 2 {
			panic(fmt.Sprintf("Line %d: function %s is invalid", lineNum+1, param))
		}
		paramType := readType(words[1], lineNum)

		parameter := variable{
			dataType: paramType,
			mutable:  false,
		}

		parameters[words[0]] = parameter
	}

	return parameters
}

func readSelectionStatement(lines []string, lineNum int, previous *selectionStatement, selectionType selectionType) (newIf selectionStatement) { // should only be called on a line after if keyword has been read
	bracketCount := 0
	var condition string

	for i := 0; i < len(lines[lineNum]); i++ {
		if lines[lineNum][i] == '(' {
			bracketCount++
		} else if lines[lineNum][i] == ')' {
			bracketCount--
		}

		if bracketCount == 1 {
			condition += string(lines[lineNum][i])
		}
	}

	// FIXME: fix massive hack below
	if expressionType([]string{removeSyntacticChars(condition)}, lineNum, nil) != Bool {
		panic(fmt.Sprintf("Line %d: if statements must have a boolean condition", lineNum+1))
	}

	newIf.condition = removeSyntacticChars(condition)
	newIf.previous = previous
	newIf.selectionType = selectionType
	newIf.begin = lineNum
	newIf.end = findScopeEnd(lines, newIf.begin)
	return newIf
}

func readSelection(lines []string, scope *scope) {
	scopeCount := 0 // used to keep track of scopes opened/closed
	for lineNum, line := range lines[(*scope).begin:(*scope).end] {
		words := strings.Fields(line)
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++
			} else if line[i] == '}' {
				scopeCount--
			}
		}
		for i := 0; i < len(words); i++ {
			if scopeCount == 1 { // 1 because if statement will have opened a scope
				if words[i] == "if" { // only read variables local to the current scope, not its subscopes.
					(*scope).selection = append((*scope).selection, readSelectionStatement(lines, (*scope).begin+lineNum, nil, If))
				} else if words[i] == "else" { // either else or else if
					if (*scope).selection[len((*scope).selection)-1].end != lineNum { // if previous selection statement does not end on line where this statement is opened
						panic(fmt.Sprintf("Line %d: else/else if statements must be opened on the same line as the corresponding if statement was closed", lineNum+1))
					}
					if words[i+1] == "if" { // else if statement
						(*scope).selection = append((*scope).selection, readSelectionStatement(lines, (*scope).begin+lineNum, &(*scope).selection[len((*scope).selection)-1], ElseIf))
					} else { // else statement
						(*scope).selection = append((*scope).selection, readSelectionStatement(lines, (*scope).begin+lineNum, &(*scope).selection[len((*scope).selection)-1], Else))
					}
					break // otherwise else ifs will trigger condition for ifs on next word
				}
			}
		}
	}
}

func readIterator(it string, lineNum int) iterator {
	newIterator := iterator{
		dataType: Int,
		start:    0,
		end:      0,
		step:     1,
	}

	parts := strings.Split(it, ":")
	fmt.Println(parts)

	if len(parts) == 2 {
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("Line %d: Start of iterators must be an integer value", lineNum))
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Sprintf("Line %d: End of iterators must be an integer value", lineNum))
		}

		if start >= end {
			panic(fmt.Sprintf("Line %d: Iterators cannot have end less than or equal to start", lineNum))
		}

		newIterator.start = start
		newIterator.end = end

	} else if len(parts) == 3 {
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("Line %d: Start of iterators must be an integer value", lineNum))
		}

		step, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Sprintf("Line %d: Step of iterators must be an integer value", lineNum))
		}

		end, err := strconv.Atoi(parts[2])
		if err != nil {
			panic(fmt.Sprintf("Line %d: End of iterators must be an integer value", lineNum))
		}

		if start >= end {
			panic(fmt.Sprintf("Line %d: Iterators cannot have end less than or equal to start", lineNum))
		}

		newIterator.start = start
		newIterator.end = end
		newIterator.step = step
	} else {
		panic(fmt.Sprintf("Line %d: Invalid iterator declaration", lineNum))
	}

	return newIterator
}

func readForLoop(lines []string, lineNum int) (newForLoop forLoop) {
	bracketCount := 0
	var it string

	for i := 0; i < len(lines[lineNum]); i++ {
		if lines[lineNum][i] == '(' {
			bracketCount++
		} else if lines[lineNum][i] == ')' && bracketCount == 1 {
			bracketCount--
			it += string(lines[lineNum][i])
		}

		if bracketCount == 1 {
			it += string(lines[lineNum][i])
		}
	}

	newForLoop.iterator = readIterator(it[1:len(it)-1], lineNum)
	newForLoop.begin = lineNum
	newForLoop.end = findScopeEnd(lines, newForLoop.begin)

	return newForLoop
}

func readIteration(lines []string, scope *scope) {
	scopeCount := 0 // used to keep track of scopes opened/closed
	for lineNum, line := range lines[(*scope).begin:(*scope).end] {
		words := strings.Fields(line)
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++
			} else if line[i] == '}' {
				scopeCount--
			}
		}
		for i := 0; i < len(words); i++ {
			if scopeCount == 1 { // 1 because if statement will have opened a scope
				if words[i] == "for" { // only read variables local to the current scope, not its subscopes.
					(*scope).iteration = append((*scope).iteration, readForLoop(lines, (*scope).begin+lineNum))
				}
			}
		}
	}
}

func findExitCondition(lines []string, begin int) { // should only be called after the loop keyword has already been read

	scopeCount := 0

	for _, line := range lines[begin:] { // first line passed in will be line where loop is opened
		words := strings.Fields(line)
		for _, word := range words {
			if word == "break" && scopeCount == 1 {
				return
			}
		}
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++
			} else if line[i] == '}' {
				scopeCount--
			}
		}
		if scopeCount == 0 {
			panic(fmt.Sprintf("Line %d: Infinite loop created without exit condition", begin))
		}
	}
}

func readInfiniteLoop(lines []string, lineNum int) (newLoop infiniteLoop) {
	findExitCondition(lines, lineNum)
	end := findScopeEnd(lines, lineNum)
	newLoop.begin = lineNum
	newLoop.end = end
	return newLoop
}

func readInfiniteLoops(lines []string, scope *scope) {
	for lineNum, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			if word == "loop" {
				(*scope).loops = append((*scope).loops, readInfiniteLoop(lines, lineNum))
			}
		}
	}
}

func readDimensions(dimensions string, lineNum int, vector bool) []int { // where e.g. "[2,3]" is passed in
	nums := strings.Split(dimensions[1:len(dimensions)-1], ",") // removes [ ] and commas
	var arrayDimensions []int
	if vector {
		for i := 0; i < len(nums); i++ {
			arrayDimensions = append(arrayDimensions, 0)
		}
	}
	for _, num := range nums {
		n, err := strconv.Atoi(num)
		if err != nil {
			panic(fmt.Sprintf("Line %d: array dimensions must be integers", lineNum))
		}
		arrayDimensions = append(arrayDimensions, n)
	}
	if len(arrayDimensions) != 0 {
		return arrayDimensions
	}
	panic(fmt.Sprintf("Line %d: arrays cannot have no specified dimensions", lineNum))
}

/*
*
func readArrayLiteral(arrayLiteral string, lineNum int) (value array) { // reads value of an array and returns value of array struct

		var currentValue string // current value being parsed
		scopeCount := 0         // important to keep track of which commas actually denote new values rather being part of the values themselves
		stringCount := 0        // same, ignores any characters which are inside a string literal
		var values []string     // slice containing values of the array
		var dimensions []int
		var currentDimensionValues int
		var elementPositions [][]int

		for i := 0; i < len(arrayLiteral); i++ { // first pass through array to figure out dimensions of array
			if arrayLiteral[i] == '"' {
				if stringCount == 0 { // if no string has already been opened
					stringCount++
				} else { // if a string has already been opened
					stringCount--
				}
			}
			if stringCount == 0 && arrayLiteral[i] == ',' {
				dimensions[scopeCount]++
			}

			if stringCount == 0 {
				if arrayLiteral[i] == '{' {
					scopeCount++
					if scopeCount > len(dimensions) {
						dimensions = append(dimensions, 0)
					}
				} else if arrayLiteral[i] == '}' {
					scopeCount--
				}
			}
		}

		for i := 0; i < len(arrayLiteral); i++ {
			if arrayLiteral[i] == '"' {
				if stringCount == 0 { // if no string has already been opened
					stringCount++
				} else { // if a string has already been opened
					stringCount--
				}
			}
			if stringCount == 0 && arrayLiteral[i] == ',' {
				dimensions[scopeCount]++
			}

			if stringCount == 0 {
				if arrayLiteral[i] == '{' {
					scopeCount++
					if scopeCount > len(dimensions) {
						dimensions = append(dimensions, 0)
					}
				} else if arrayLiteral[i] == '}' {
					scopeCount--
				}
			}

		}

		// TODO: next pass through array actually assigning values to the dimensions that have been read
		// TODO: do this by creating values slice based on dimensions slice

		value.dataType = getValType(values[0], lineNum)
		for _, val := range values {
			if getValType(val, lineNum) != value.dataType {
				panic(fmt.Sprintf("Line %d: all values in an array must be of the same type", lineNum))
			}
		}

		return value
	}
*/
func readArray(lines []string, lineNum int) (newArray array) { // only called after arr or vec keyword detected
	words := strings.Fields(lines[lineNum])
	if words[0] == "vec" {
		newArray.vec = true // false by default
	}

	var arrayType string
	var arrayDimensions string
	// var arrayLiteral string

	for i := 0; i < len(words); i++ {
		if words[i][len(words[i])-1] == ':' { // last character of colon indicates that this is the variable name and the type is next
			scopeCount := 0
			for j := 0; j < len(words[i+1]); j++ {
				if words[i+1][j] == '[' {
					scopeCount++
				}
				if scopeCount == 0 {
					arrayType += string(words[i+1][j])
				} else if scopeCount == 1 {
					arrayDimensions += string(words[i+1][j])
				}

			}
			newArray.dataType = readType(words[i+1], lineNum)

		} else if words[i] == "=" {
			if words[i+1][0] != '{' || words[i+1][len(words[i+1-1])] != '}' {
				panic(fmt.Sprintf("Line %d: array literals must be enclosed by curly brackets", lineNum))
			} else {
				// arrayLiteral = words[i+1]
			}
		}
	}

	if !newArray.vec {
		newArray.dimensions = readDimensions(arrayDimensions, lineNum, false)
	} else {
		newArray.dimensions = readDimensions(arrayDimensions, lineNum, true)
	}

	newArray.dataType = readType(arrayType, lineNum)
	//	readArrayLiteral(arrayLiteral, lineNum)

	return newArray
}

func readArrays(lines []string, scope *scope) {
	scopeCount := 0 // used to keep track of scopes opened/closed
	for lineNum, line := range lines[(*scope).begin:(*scope).end] {
		words := strings.Fields(line)
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++
			} else if line[i] == '}' {
				scopeCount--
			}
		}
		for i := 0; i < len(words); i++ {
			if scopeCount == 1 { // 1 because scopeCount will be 1 from line 1
				if words[i] == "arr" || words[i] == "vec" { // only read variables local to the current scope, not its subscopes.
					(*scope).arrs[readName(lines, lineNum)] = readArray(lines, lineNum)
				}
			}
		}
	}
}

/**
func readAssignment(lines []string, lineNum int, currentScope *scope) {
	line := lines[lineNum]
	name := strings.Fields(line)[0]
	if variable, ok := currentScope.vars[name]; ok {
		if !variable.mutable {
			(*currentScope).vars[name] = variable
		}
		panic(fmt.Sprintf("Line %d: Cannot mutate constant %s", lineNum+1, name))
	} else { // variable name not found
		panic(fmt.Sprintf("Line %d: variable name %s does not exist", lineNum+1, name))
	}
}
*/
