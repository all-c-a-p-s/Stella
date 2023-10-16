package main

import (
	"fmt"
	"strings"
)

type primitiveType int

const (
	String primitiveType = iota
	Int
	Bool
	Float
	Invalid
)

func readType(dataType string, lineNum int) primitiveType { //reads data type from assignment
	switch dataType {
	case "int":
		return Int
	case "float":
		return Float
	case "bool":
		return Bool
	case "string":
		return String
	default:
		panic(fmt.Sprintf("Line %d - data type %s is invalid", lineNum+1, dataType))
	}
}

func getValType(value string, lineNum int) primitiveType {
	if value[0] == '"' && value[len(value)-1] == '"' {
		checkStringVal(value, lineNum)
		return String
	} else if value == "true" || value == "false" {
		checkBoolVal(value, lineNum)
		return Bool
	} // else must be a number

	for _, char := range value {
		if char == '.' {
			checkFloatVal(value, lineNum)
			return Float
		}
	}

	checkIntVal(value, lineNum)
	return Int
}

func getType(expression string, lineNum int, currentScope *scope) primitiveType {
	//first check if expression contains variable names
	//variable names will be either their own word, or next to annotation characters
	//TODO: implement this for function names as well
	words := strings.Fields(expression)
	typesFound := []primitiveType{}
	for _, word := range words {
		if variable, ok := (*currentScope).vars[removeSyntacticChars(word)]; ok {
			typesFound = append(typesFound, variable.dataType)
		} else { //not a variable name so it must be a value
			typesFound = append(typesFound, getValType(removeSyntacticChars(word), lineNum))
		}
	}
	for i := 1; i < len(typesFound); i++ { //checks that all types found are the same
		if typesFound[i] != typesFound[0] {
			panic(fmt.Sprintf("Line %d - Expression contains different types", lineNum+1))
		}
	}

	return typesFound[0] //all types in expression are the same, and the expression must have a type so this works
}

func checkIntVal(value string, lineNum int) { //checks to see if int value contains illegal characters/leading zeros etc.
	switch value[0] {
	case '0':
		panic(fmt.Sprintf("Line %d - Integers values cannot have leading zeros", lineNum+1))
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		panic(fmt.Sprintf("Line %d - character %c cannot be part of an integer value", lineNum+1, value[0]))
	}

	for _, char := range value {
		if !(char > 47 && char < 58) { //digits including zero. leading zeros will have been caught above
			panic(fmt.Sprintf("Line %d - character %c cannot be part of an integer value", lineNum+1, char))
		}
	}

}

func checkFloatVal(value string, lineNum int) {
	switch value[0] {
	case '0':
		if !(value[1] == '.') {
			panic(fmt.Sprintf("Line %d - Leading zeros must be followed by decimal point, here it is followed by %c", lineNum+1, value[1]))
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		panic(fmt.Sprintf("Line %d - character %c cannot be part of an float value", lineNum+1, value[0]))
	}

	decimalPointCount := 0
	for _, char := range value {
		if char == '.' {
			decimalPointCount++
		}
		if !(char > 47 && char < 58) { //digits including zero. leading zeros will have been caught above
			if !(char == '.' && decimalPointCount == 0) {
				panic(fmt.Sprintf("Line %d - character %c cannot be part of an integer value", lineNum+1, char))
			}
		}
	}
}

func checkBoolVal(value string, lineNum int) {
	if !(value == "true" || value == "false") {
		panic(fmt.Sprintf("Line %d - value '%s' cannot be used as a boolean value", lineNum+1, value))
	}
}

func checkStringVal(value string, lineNum int) {
	if !(value[0] == '"' && value[len(value)-1] == '"') {
		panic(fmt.Sprintf("Line %d - '%s' cannot be used as string value", lineNum+1, value))
	}
}
