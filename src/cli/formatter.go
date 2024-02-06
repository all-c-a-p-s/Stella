package main

import "strings"

// TODO: check (somehow) that this works on all OS, architecture
func format(transpiled string, tabSize int) string {
	lines := strings.Split(transpiled, "\n")
	var bracketCount int
	var bracketScoreAtEnd int
	var stringLiteralCount int

	var formatted string

	for _, line := range lines {
		for spaces := 0; spaces < bracketScoreAtEnd; spaces++ {
			for k := 0; k < tabSize; k++ {
				formatted += " "
			}
		}
		formatted += line
		for i := 0; i < len(line); i++ {
			switch line[i] {
			case '"':
				if stringLiteralCount == 1 {
					stringLiteralCount = 0
				} else {
					stringLiteralCount = 1
				}
			case '{':
				if stringLiteralCount == 0 {
					bracketCount++
				}
			case '}':
				if stringLiteralCount == 0 {
					bracketCount--
				}
			}
		}
		formatted += "\n"
		bracketScoreAtEnd = bracketCount
	}
	return formatted
}
