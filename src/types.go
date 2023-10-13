package main

import (
	"fmt"
)

type primitiveType int

const (
	String primitiveType = iota
	Int
	Bool
	Float
)

func getType(value string, lineNum int) primitiveType {
	if value[0] == '"' && value[len(value)-1] == '"' {
		checkString(value, lineNum)
		return String
	} else if value == "true" || value == "false" {
		checkBool(value, lineNum)
		return Bool
	} // else must be a number

	for _, char := range value {
		if char == '.' {
			checkFloat(value, lineNum)
			return Float
		}
	}

	checkInt(value, lineNum)
	return Int
}

func checkInt(value string, lineNum int) { //checks to see if int value contains illegal characters/leading zeros etc.
	switch value[0] {
	case '0':
		panic(fmt.Sprintf("Line %d - Integers values cannot have leading zeros", lineNum))
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		panic(fmt.Sprintf("Line %d - character %c cannot be part of an integer value", lineNum, value[0]))
	}

	for _, char := range value {
		if !(char > 47 && char < 58) { //digits including zero. leading zeros will have been caught above
			panic(fmt.Sprintf("Line %d - character %c cannot be part of an integer value", lineNum, char))
		}
	}

}

func checkFloat(value string, lineNum int) {
	switch value[0] {
	case '0':
		if !(value[1] == '.') {
			panic(fmt.Sprintf("Line %d - Leading zeros must be followed by decimal point, here it is followed by %c", lineNum, value[1]))
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		panic(fmt.Sprintf("Line %d - character %c cannot be part of an float value", lineNum, value[0]))
	}

	decimalPointCount := 0
	for _, char := range value {
		if char == '.' {
			decimalPointCount++
		}
		if !(char > 47 && char < 58) { //digits including zero. leading zeros will have been caught above
			if !(char == '.' && decimalPointCount == 0) {
				panic(fmt.Sprintf("Line %d - character %c cannot be part of an integer value", lineNum, char))
			}
		}
	}
}

func checkBool(value string, lineNum int) {
	if !(value == "true" || value == "false") {
		panic(fmt.Sprintf("Line %d - value '%s' cannot be used as a boolean value", lineNum, value))
	}
}

func checkString(value string, lineNum int) {
	if !(value[0] == '"' && value[len(value)-1] == '"') {
		panic(fmt.Sprintf("Line %d - '%s' cannot be used as string value", lineNum, value))
	}
}
