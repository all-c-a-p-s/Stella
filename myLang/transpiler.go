package main 

import (
	"strings"	
	"os"
)

type scope struct {
	begin int
	end int
	subScopes []*scope 
}

func check(err error) {
	if err != nil {
		panic(err)		
	}
}

//here lineIndex is the start of the scope
func findScopeEnd(lines []string, lineIndex int) int {
	for lineNum, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			if word == "}" {
				return lineIndex + lineNum
			}
		}
	}
	panic("Could not find end of scope")
}


//string is everyting inside current scope
//begin and end are line numbers of beginning and end of current scope
func findSubScopes(lines []string, begin, end int, currentScope *scope) []*scope {

	//use recursion to find sub-scopes

	for lineNum, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			if word == "{" { //denotes beginning of a scope
				scopeBeginning := begin + lineNum
				scopeEnd := findScopeEnd(lines[begin + lineNum:end], begin + lineNum)

				var subScope scope //initialise so function call can have pointer to itself

				subScope.begin = scopeBeginning
				subScope.end = scopeEnd
				subScope.subScopes = findSubScopes(lines[scopeBeginning:scopeEnd], scopeBeginning, scopeEnd, &subScope) //updates sub-scope list of sub-scope
				
				(*currentScope).subScopes = append((*currentScope).subScopes, &subScope)
			}
		}			
	}

	return []*scope{} //returns empty slice if none found
}

func main() {	
	src, err := os.ReadFile("src.txt")
	check(err)

	sourceCode := string(src)
	lines := strings.Split(sourceCode, "\n")

	var globalScope scope

	globalScope.begin = 0
	globalScope.end = len(lines)
	globalScope.subScopes = findSubScopes(lines, 0, len(lines), &globalScope)
}
