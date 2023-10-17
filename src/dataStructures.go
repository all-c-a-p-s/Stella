package main

import (
	"fmt"
	"strings"
)

type variable struct { //name not used here because it will be stored as the key of the hashmap of variables
	dataType primitiveType
	constant bool
}

type function struct {
	returnType primitiveType
	parameters map[string]variable
}

type expression struct {
	dataType string
}

func primitiveTypes() []primitiveType {
	return []primitiveType{Int, Float, Bool, String}
}

func charType(char byte) string {
	switch char {
	case 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90: //uppercase letter
		return "letter"
	case 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122: //lowercase
		return "letter"
	case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
		return "number"
	case 95:
		return "underscore"
	default:
		return "other"
	}

}

func validName(name string, lineNum int) string {
	if !(charType(name[0]) == "letter") { //doesn't begin with uppercase or lowercase letter
		panic(fmt.Sprintf("Line %d - Name '%s' is invalid because it does not begin with a letter", lineNum+1, name))
	}

	for i := 0; i < len(name)-1; i++ { //last character can be syntactic character
		if !(charType(name[i]) == "letter" || charType(name[i]) == "number" || charType(name[i]) == "underscore") { //character other than letters, number or underscore
			panic(fmt.Sprintf("Line %d - Name '%s' is invalid because it contains invalid character '%s'", lineNum+1, name, string(name[i])))
		}
	}

	last := len(name) - 1

	if !(charType(name[last]) == "letter" || charType(name[last]) == "number" || syntacticCharacter(name[last])) { //last character can be annotation character
		panic(fmt.Sprintf("Line %d - Name '%s' is invalid because it contains invalid last character '%s'", lineNum+1, name, string(name[len(name)-1])))
	}
	return name //no exit conditions triggered, so name must be valid
}

func syntacticCharacter(char byte) bool { //characters which have a syntactic function
	switch char {
	case ':', '(', ')', '[', ']', '{', '}':
		return true
	default:
		return false
	}
}

func declarationKeyword(word string) bool {
	if word == "let" || word == "const" || word == "func" {
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
			if syntacticCharacter(words[i+1][len(words[i+1])-1]) { //names can have syntactic characters e.g. ':' or '(' after them without a space
				return validName(words[i+1][:len(words[i+1])], lineNum)
			}
			return validName(words[i+1], lineNum)
		} else if words[i] == "func" {
			var funcName string
			for j := 0; j < len(words[i+1]); j++ {
				if words[i+1][j] == '(' { //start of function parameters
					break
				}
				funcName += string(words[i+1][j])
			}
		}
	}
	panic(fmt.Sprintf("Line %d - invalid name in declaration", lineNum+1))
}

func readVariable(lines []string, lineNum int) (newVariable variable) {

	var dataType primitiveType
	constant := false

	words := strings.Fields(lines[lineNum])

	if words[0] == "const" {
		constant = true
	}

	for i := 0; i < len(words); i++ {
		if words[i][len(words[i])-1] == ':' { //last character of colon indicates that this is the variable name and the type is next
			dataType = readType(words[i+1], lineNum)

		}
	}
	newVariable.dataType = dataType
	newVariable.constant = constant
	return newVariable
}

func readVariables(lines []string, scope *scope) {
	scopeCount := 0 //used to keep track of scopes opened/closed
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
			if word == "let" && scopeCount == 0 { //only read variables local to the current scope, not its subscopes
				(*scope).vars[readName(lines, lineNum)] = readVariable(lines, lineNum)
			}
		}
	}
}

func readFunction(lines []string, lineNum int) (newFunction function) {
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
		returnType = Void
	}

	newFunction.returnType = returnType
	newFunction.parameters = readParameters(words[1], lineNum)

	return newFunction
}

func readFunctions(lines []string, scope *scope) {
	scopeCount := 0 //used to keep track of scopes opened/closed
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
			if word == "func" && scopeCount == 0 { //only read variables local to the current scope, not its subscopes
				(*scope).functions[readName(lines, lineNum)] = readFunction(lines, lineNum)
			}
		}
	}
}

func readParameters(funcName string, lineNum int) (parameters map[string]variable) {
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

	params = removeSyntacticChars(params) //removes opening bracket
	funcParams := strings.Split(params, ",")

	for _, param := range funcParams {
		words := strings.Fields(param)
		if len(words) != 2 {
			panic(fmt.Sprintf("Line %d - function %s is invalid", lineNum+1, param))
		}
		paramType := readType(words[1], lineNum)

		parameter := variable{
			dataType: paramType,
			constant: true,
		}

		parameters[words[0]] = parameter
	}

	return parameters
}

func assignValue(lines []string, lineNum int, currentScope *scope) {
	line := lines[lineNum]
	name := strings.Fields(line)[0]
	if variable, ok := currentScope.vars[name]; ok {
		if !variable.constant {
			(*currentScope).vars[name] = variable
		}
		panic(fmt.Sprintf("Line %d - Cannot mutate constant %s", lineNum+1, name))
	} else { //variable name not found
		panic(fmt.Sprintf("Line %d - variable name %s does not exist", lineNum+1, name))
	}
}
