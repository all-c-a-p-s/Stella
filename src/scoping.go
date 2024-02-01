package main

import (
	"bufio"
	"fmt"
	"os"
)

type ScopeType int

const (
	FunctionScope ScopeType = iota
	SelectionScope
	Global
)

type Scope struct {
	vars      map[string]Variable
	functions map[string]Function
	items     []Transpileable
	scopeType ScopeType
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

	globalScope := parseScope(lines, 0, Global, nil)
	fmt.Println(parseArray("[1, 2, 3]", Int, &globalScope, 0))
	fmt.Println(globalScope.transpile())
}
