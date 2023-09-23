package main

import (
	"bufio"
	"os"
)

type scope struct {
	begin     int
	end       int
	subScopes []*scope
	parent    *scope
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// here lineIndex is the start of the scope
func findScopeEnd(lines []string, lineIndex int) int {

	//scopeCount keeps track of scopes opened/closed. When this is zero the current scope has been closed
	scopeCount := 0 //initialise at 0 because it will immediately become 1 after it reads first line

	for lineNum, line := range lines {

		for i := 0; i < len(line); i++ {
			if line[i] == '{' { //new scope opened
				scopeCount++
			} else if line[i] == '}' { //scope closed
				scopeCount--
			}
		}

		if scopeCount == 0 {
			return lineIndex + lineNum
		}
	}
	panic("Could not find end of scope")
}

// lines inside scope passed in as slice all lines in the file
// begin and end are line numbers of beginning and end of current scope
func findSubScopes(lines []string, begin, end int, currentScope *scope) []*scope {

	subScopes := []*scope{}
	scopeCount := 0 //keeps track of scopes opened/closed. scopeCount of 1 indicates a subScope of the current scope

	for lineNum, line := range lines {
		for i := 0; i < len(line); i++ { // better than using strings.Fields() in case user doesn't put a space before opening a scope
			if line[i] == '{' {
				scopeCount++
				if scopeCount == 1 { //should only be exectued on lines where a scope is actually opened
					scopeBeginning := begin + lineNum
					scopeEnd := findScopeEnd(lines[scopeBeginning:end], scopeBeginning)

					var subScope scope

					subScope.begin = scopeBeginning
					subScope.end = scopeEnd
					subScope.subScopes = []*scope{}
					subScope.parent = currentScope

					//recursive call to findSubScopes() to generate tree of scopes/subScopes
					subScope.subScopes = findSubScopes(lines[scopeBeginning+1:end], scopeBeginning, scopeEnd, &subScope)

					subScopes = append(subScopes, &subScope)
				}
			} else if line[i] == '}' {
				scopeCount--
			}
		}

	}

	return subScopes
}

func main() {
	src, err := os.Open("src.txt")
	check(err)

	//scanner used to avoid OS-specific problems, e.g. windows having "\r\n" for newlines rather than just "\n"
	scanner := bufio.NewScanner(src)

	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text()) //read text from each lines slice
	}

	//globalScope is the scope of the entire program (i.e. everything)
	var globalScope scope

	globalScope.begin = 0
	globalScope.end = len(lines)
	globalScope.subScopes = []*scope{}
	globalScope.parent = nil

	globalScope.subScopes = findSubScopes(lines, 0, len(lines), &globalScope)
}
