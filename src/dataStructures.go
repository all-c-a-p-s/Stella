package main

import (
	"fmt"
	"strings"
)

type selectionType int

const (
	If selectionType = iota
	ElseIf
	Else
)

type variable struct { //name not used here because it will be stored as the key of the hashmap of variables
	dataType primitiveType
	constant bool
}

type function struct {
	returnType primitiveType
	parameters map[string]variable
}

type selectionStatement struct {
	previous      *selectionStatement //previous if statement e.g. if statement before else if
	condition     string
	selectionType selectionType
	begin         int //first line
	end           int //last line of statement, these are used to check if else if/else statements are opened on the same line as if statements are closed
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

func booleanOperator(operator string) bool {
	switch operator {
	case "==", "<=", ">=", "<", ">", "!=":
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

func readFunction(lines []string, lineNum int) (newFunction function) { //should only be called on a line once the func keyword has already been read
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

func booleanExpression(expression string) bool { //only searches for boolean operator, does not check if expression is valid
	words := strings.Fields(expression)
	for _, word := range words {
		if booleanOperator(word) {
			return true
		}
	}

	return false
}

func readSelectionStatement(lines []string, lineNum int, previous *selectionStatement, selectionType selectionType) (newIf selectionStatement) { //should only be called on a line after if keyword has been read
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

	if !(booleanExpression(removeSyntacticChars(condition))) {
		panic(fmt.Sprintf("Line %d - if statements must have a boolean condition", lineNum+1))
	}

	newIf.condition = removeSyntacticChars(condition)
	newIf.previous = previous
	newIf.selectionType = selectionType
	newIf.begin = lineNum
	newIf.end = findScopeEnd(lines, newIf.begin)
	return newIf
}

func readSelection(lines []string, scope *scope) {
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
		for i := 0; i < len(words); i++ {
			if scopeCount == 1 { //1 because if statement will have opened a scope
				if words[i] == "if" { //only read variables local to the current scope, not its subscopes.
					(*scope).selection = append((*scope).selection, readSelectionStatement(lines, (*scope).begin+lineNum, nil, If))
				} else if words[i] == "else" { //either else or else if
					if (*scope).selection[len((*scope).selection)-1].end != lineNum { //if previous selection statement does not end on line where this statement is opened
						panic(fmt.Sprintf("Line %d - else/else if statements must be opened on the same line as the corresponding if statement was closed", lineNum+1))
					}
					if words[i+1] == "if" { //else if statement
						(*scope).selection = append((*scope).selection, readSelectionStatement(lines, (*scope).begin+lineNum, &(*scope).selection[len((*scope).selection)-1], ElseIf))
					} else { //else statement
						(*scope).selection = append((*scope).selection, readSelectionStatement(lines, (*scope).begin+lineNum, &(*scope).selection[len((*scope).selection)-1], Else))
					}
					break //otherwise else ifs will trigger condition for ifs on next word
				}
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
		panic(fmt.Sprintf("Line %d - Cannot mutate constant %s", lineNum+1, name))
	} else { //variable name not found
		panic(fmt.Sprintf("Line %d - variable name %s does not exist", lineNum+1, name))
	}
}
