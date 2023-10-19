package main

import (
	"bufio"
	"fmt"
	"os"
)

type scope struct {
	begin     int //0-indexed, inclusive
	end       int //0-indexed, exclusive
	subScopes []*scope
	parent    *scope
	vars      map[string]variable
	functions map[string]function
	selection []selectionStatement
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func findScopeEnd(lines []string, begin int) int {

	scopeCount := 0 //keeps track of scopes opened/scopes closed

	for lineNum, line := range lines[begin:] { //first line passed in will be line where scope is opened
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++
			} else if line[i] == '}' && lineNum+begin != begin { //important for lines where a scope is opened on the same line where another is closed
				scopeCount-- //for lines which do not close the current scope but do close a scope
				if scopeCount == 0 {
					return begin + lineNum
				}
			}
		}
	}

	panic(fmt.Sprintf("Line %d - scope opened but never closed", begin+1))
}

func readScope(lines []string, begin, end int, currentScope *scope) {

	readVariables(lines, currentScope)
	readFunctions(lines, currentScope)
	readSelection(lines, currentScope)

	scopeCount := 0 //keeps track of scopes opened/scopes closed. Count of 2 will indicate a new subscope being opened

	if (*currentScope).parent == nil { //global scope
		scopeCount++ //incremented because the global scope is the only scope where a bracket is not used to open the scope
	}

	for lineNum, line := range lines[begin:end] {
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++

				if scopeCount == 2 { //only execute on lines where a scope is actually opened
					scopeBeginning := begin + lineNum
					scopeEnd := findScopeEnd(lines, scopeBeginning) //lines[lineNum:] because the slice of the function parameter is passed
					//to findScopeEnd, so we need the relative position

					subScope := scope{
						begin:     scopeBeginning,
						end:       scopeEnd,
						subScopes: []*scope{},
						parent:    currentScope,
						vars:      map[string]variable{},
					}

					(*currentScope).subScopes = append((*currentScope).subScopes, &subScope)
					readScope(lines, scopeBeginning, scopeEnd, &subScope)
				}

			} else if line[i] == '}' {
				scopeCount--
			}
		}
	}
}

func main() {
	src, err := os.Open("src.txt")
	check(err)
	defer src.Close()

	var lines []string //all lines of source code will be passed into functions

	scanner := bufio.NewScanner(src) //used to avoid OS-specific problems such as Windows using "\r\n" for newline rather than just "\n"

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	globalScope := scope{ //where globalScope is the entire file
		begin:     0,
		end:       len(lines),
		subScopes: []*scope{},
		parent:    nil,
		vars:      map[string]variable{},
	}

	readScope(lines, 0, len(lines), &globalScope)
	fmt.Println("Compiled successfully")
	fmt.Println(globalScope)
}
