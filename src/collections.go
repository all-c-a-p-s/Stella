package main

import (
	"fmt"
	"strconv"
)

// let nums: int[5] = [1, 2, 3, 4, 5]
type ArrayType struct {
	dimensions []int
	baseType   primitiveType
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
