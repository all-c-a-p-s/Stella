package main

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
	arr  Array
	expr ArrayValue[primitiveType]
}

func squareBracketEnd(s string, start, lineNum int) int {
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

	dims := typeWord[squareBracketIndex:]
	bracketCount := 0

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

	// TODO: fix string literals

	for i := 1; i < len(arrayValue); i++ {
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

func parseArray[T primitiveType](arrayValue string, expectedType primitiveType, currentScope *Scope, lineNum int) ArrayValue[T] {
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
				child := parseArray[T](subArrayValue, expectedType, currentScope, lineNum)
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
	arrFound := parseArray[primitiveType](expression, expectedType.baseType, currentScope, lineNum)

	if len(arrFound.children) > 0 {
		panic(fmt.Sprintf("Line %d: currently elements of arrays can only be primitive types, not arrays", lineNum+1))
	}

	if arrFound.baseType != expectedType.baseType {
		panic(fmt.Sprintf("Line %d: expected array of type %v found array of type %v", lineNum+1, expectedType.baseType, arrFound.baseType))
	}

	if arrFound.length != expectedType.dimensions[0] {
		panic(fmt.Sprintf("Line %d: expected array of length %d, found array of length %d", lineNum+1, arrFound.length, expectedType.dimensions[0]))
	}

	arr := Array{
		mut:        mut,
		identifier: id,
		dataType:   expectedType,
	}
	return ArrayDeclaration{
		arr:  arr,
		expr: arrFound,
	}
}
