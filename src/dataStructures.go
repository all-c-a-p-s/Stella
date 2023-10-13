package main

import (
	"fmt"
	"slices"
	"strings"
)

type variable struct { //name not used here because it will be stored as the key of the hashmap of variables
	dataType  string
	initScope *scope
	constant  bool
}

type expression struct {
	dataType string
}

func primitiveTypes() []string {
	return []string{"int", "float", "bool", "string"}
}

func validName(name string, lineNum int) string {
	if !((name[0] >= 65 && name[0] <= 90) || (name[0] >= 97 && name[0] <= 122)) { //doesn't begin with uppercase or lowercase letter
		panic(fmt.Sprintf("Line %d - Name '%s' is invalid because it does not begin with a letter", lineNum+1, name))
	}

	for i := 0; i < len(name); i++ {
		if !((name[i] >= 65 && name[i] <= 90) || (name[i] >= 97 && name[i] <= 122) || (name[i] == 95)) { //character other than letters or underscore
			panic(fmt.Sprintf("Line %d - Name '%s' is invalid because it contains invalid character '%s'", lineNum+1, name, string(name[i])))
		}
	}
	return name //no exit conditions triggered, so name must be valid
}

func annotationCharacter(char byte) bool { //characters which have a syntactic function
	switch char {
	case ':' | '(' | ')' | '[' | ']' | '{' | '}':
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

func readName(lines []string, lineNum int) string {
	line := lines[lineNum]
	words := strings.Fields(line)

	for i := 0; i < len(words); i++ {
		if declarationKeyword(words[i]) {
			if annotationCharacter(words[i+1][len(words[i+1])-1]) { //names can have syntactic characters e.g. ':' or '(' after them without a space
				return validName(words[i+1][:len(words[i+1])], lineNum)
			}
			return validName(words[i+1], lineNum)
		}
	}
	panic(fmt.Sprintf("Line %d - invalid name in declaration", lineNum+1))
}

func readVariable(lines []string, lineNum int) (newVariable variable) {

	var name string
	var dataType string
	constant := false

	words := strings.Fields(lines[lineNum])

	if words[0] == "const" {
		constant = true
	}

	for i := 0; i < len(words); i++ {
		if words[i][len(words[i])-1] == ':' { //last character of colon indicates that this is the variable name and the type is next
			name = words[i][:len(words[i])-1] //name without the colon
			if slices.Contains(primitiveTypes(), words[i+1]) {
				dataType = words[i+1]
			}
			panic(fmt.Sprintf("Line %d - Variable %s has invalid type", lineNum+1, name))
		}
		panic(fmt.Sprintf("Line %d - Variable declaration without type annotation", lineNum+1))
	}
	newVariable.dataType = dataType
	newVariable.constant = constant
	panic("Keyword 'let' used without variable declaration")
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

func assignValue(lines []string, lineNum int, currentScope *scope) {
	line := lines[lineNum]
	name := strings.Fields(line)[0]
	if variable, ok := currentScope.vars[name]; ok {
		if !variable.constant {
			(*currentScope).vars[name] = variable
		}
		panic(fmt.Sprintf("Line %d - Cannot mutate constant %s", lineNum, name))
	}

	panic(fmt.Sprintf("Line %d - variable name %s does not exist", lineNum, name))
}

func readExpression(lines []string, lineNum int) {

}
