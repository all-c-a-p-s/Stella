package main

import (
	"bufio"
	"os"
	"slices"
	"strings"
)

type variable struct {
	name     string
	dataType string
	value    string
}

type scope struct {
	begin     int //0-indexed, inclusive
	end       int //0-indexed, exclusive
	subScopes []*scope
	parent    *scope
	vars      []variable
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func findScopeEnd(lines []string, startLine int) int {

	scopeCount := 0 //keeps track of scopes opened/scopes closed

	for lineNum, line := range lines { //first line passed in will be line where scope is opened
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++
			} else if line[i] == '}' {
				scopeCount--
			}
		}

		if scopeCount == 0 {
			return startLine + lineNum
		}
	}

	panic("Could not find end of scope")
}

func readScope(lines []string, startLine int, currentScope *scope) {

	if (*currentScope).parent != nil { //not for global scope which has nil parent address
		(*currentScope).vars = (*currentScope).parent.vars //inherits variables from parent scope
	}
	readVariables(lines, currentScope)

	scopeCount := 0 //keeps track of scopes opened/scopes closed. Count of 2 will indicate a new subscope being opened

	if (*currentScope).parent == nil { //global scope
		scopeCount++ //incremented because the global scope is the only scope where a bracket is not used to open the scope
	}

	for lineNum, line := range lines {
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++

				if scopeCount == 2 { //only execute on lines where a scope is actually opened
					scopeBeginning := startLine + lineNum
					scopeEnd := findScopeEnd(lines[lineNum:], scopeBeginning) //lines[lineNum:] because the slice of the function parameter is passed
					//to findScopeEnd, so we need the relative position

					subScope := scope{
						begin:     scopeBeginning,
						end:       scopeEnd,
						subScopes: []*scope{},
						parent:    currentScope,
						vars:      []variable{},
					}

					readScope(lines[scopeBeginning:scopeEnd], scopeBeginning, &subScope)

					(*currentScope).subScopes = append((*currentScope).subScopes, &subScope)
				}

			} else if line[i] == '}' {
				scopeCount--
			}
		}
	}
}

func readValue(line string) string {
	words := strings.Fields(line)

	for i := 0; i < len(words); i++ {
		if words[i] == "=" {
			return words[i+1]
		}
	}
	panic("Failed to read value")
}

func readVariable(line string) variable {
	var types = []string{"int", "string", "bool", "char", "float"}

	words := strings.Fields(line)
	for i := 0; i < len(words); i++ {
		if words[i][len(words[i])-1] == ':' { //last character of colon indicates that this is the variable name and the type is next
			if slices.Contains(types, words[i+1]) {
				newVariable := variable{
					name:     words[i][:len(words[i])-1], //name without the colon
					dataType: words[i+1],
					value:    readValue(line),
				}
				return newVariable
			}
		}
	}

	panic("Failed to read variable declaration")
}

func readVariables(lines []string, scope *scope) {
	scopeCount := 0 //used to keep track of scopes opened/closed
	for _, line := range lines {
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
				(*scope).vars = append((*scope).vars, readVariable(line))
			}
		}
	}
}

func main() {
	src, err := os.Open("src.txt")
	check(err)
	defer src.Close()

	var lines []string

	scanner := bufio.NewScanner(src) //used to avoid OS-specific problems such as Windows using "\r\n" for newline rather than just "\n"

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	globalScope := scope{ //where globalScope is the entire file
		begin:     0,
		end:       len(lines),
		subScopes: []*scope{},
		parent:    nil,
		vars:      []variable{},
	}

	readScope(lines, 0, &globalScope)
}
