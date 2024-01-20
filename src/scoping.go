package main

import (
	"bufio"
	"fmt"
	"os"
)

type scope struct {
	begin     int // 0-indexed, inclusive
	end       int // 0-indexed, exclusive
	subScopes []*scope
	parent    *scope
	vars      map[string]Variable
	functions map[string]function
	arrs      map[string]array
	selection []selectionStatement
	iteration []forLoop
	loops     []infiniteLoop
}

type Location struct {
	// location in source file
	lineNum   int
	charIndex int
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func findScopeEnd(lines []string, begin int) int {
	scopeCount := 0 // keeps track of scopes opened/scopes closed
	opened := false // keeps track of if scope has been opened yet. important for lines where a scope if opened on the same line where another is closed

	for lineNum, line := range lines[begin:] { // first line passed in will be line where scope is opened
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				opened = true
				scopeCount++
			} else if line[i] == '}' && opened {
				scopeCount-- // for lines which do not close the current scope but do close a scope
				if scopeCount == 0 {
					return begin + lineNum
				}
			}
		}
	}

	panic(fmt.Sprintf("Line %d: scope opened but never closed", begin+1))
}

func findBracketEnd(bracketType byte, lines []string, lineNum int, charIndex int) Location {
	// should be called wher lineNum and charIndex are the location of the character operning the brackets
	// this means that it will become 1 on the first character
	bracketCount := 0
	var closingBracket byte
	switch bracketType {
	case '(':
		closingBracket = ')'
	case '{':
		closingBracket = '}'
	case '[':
		closingBracket = ']'
	default:
		panic("Invalid character used as bracketType")
	}
	for i := lineNum; i < len(lines); i++ {
		line := lines[i]
		var charStart int // character to start searching the line
		if i == lineNum { // only on start line
			charStart = charIndex
		}
		for j := charStart; j < len(line); j++ {
			switch line[j] {
			case bracketType:
				bracketCount++
			case closingBracket:
				bracketCount--
			}
			if bracketCount == 0 { // 0 at end of loop means it must have been closed
				return Location{lineNum: i, charIndex: j}
			}
		}

	}
	panic(fmt.Sprintf("Line %d: bracket %s opened but never closed", lineNum+1, string(bracketType)))
}

func readScope(lines []string, begin, end int, currentScope *scope) {
	// readVariables(lines, currentScope)
	readFunctions(lines, currentScope)
	readSelection(lines, currentScope)
	readIteration(lines, currentScope)
	readInfiniteLoops(lines, currentScope)

	scopeCount := 0 // keeps track of scopes opened/scopes closed. Count of 2 will indicate a new subscope being opened

	if (*currentScope).parent == nil { // global scope
		scopeCount++ // incremented because the global scope is the only scope where a bracket is not used to open the scope
	}

	for lineNum, line := range lines[begin:end] {
		for i := 0; i < len(line); i++ {
			if line[i] == '{' {
				scopeCount++

				if scopeCount == 2 { // only execute on lines where a scope is actually opened
					scopeBeginning := begin + lineNum
					scopeEnd := findScopeEnd(lines, scopeBeginning) // lines[lineNum:] because the slice of the function parameter is passed
					// to findScopeEnd, so we need the relative position

					subScope := scope{
						begin:     scopeBeginning,
						end:       scopeEnd,
						subScopes: []*scope{},
						parent:    currentScope,
						vars:      map[string]Variable{},
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
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			panic("Error closing source code file")
		}
	}(src)

	var lines []string // all lines of source code will be passed into functions

	scanner := bufio.NewScanner(src) // used to avoid OS-specific problems such as Windows using "\r\n" for newline rather than just "\n"

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	globalScope := scope{ // where globalScope is the entire file
		begin:     0,
		end:       len(lines),
		subScopes: []*scope{},
		parent:    nil,
		vars:      map[string]Variable{},
	}

	var fn string = "fn square(x: int) -> int = { \n 5 * 5 \n }"
	lns := []string{fn}
	fmt.Println(parseFunction(lns, 0, &globalScope))
	// readScope(lines, 0, len(lines), &globalScope)
	// fmt.Println("Compiled successfully")
	// fmt.Println(globalScope.subScopes[0])
}
