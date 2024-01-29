package main

func binaryOperators() map[string]struct{} {
	// operators that need an expression on both sides
	return map[string]struct{}{
		"+":  {},
		"-":  {},
		"*":  {},
		"/":  {},
		"&&": {},
		"||": {},
		"==": {},
		"!=": {},
		">":  {},
		"<":  {},
		"<=": {},
		">=": {},
	}
}

func numericOperators() map[string]struct{} {
	return map[string]struct{}{
		"+": {},
		"-": {},
		"*": {},
		"/": {},
	}
}

func comparativeOperators() map[string]struct{} {
	return map[string]struct{}{
		"==": {},
		"!=": {},
		">":  {},
		"<":  {},
		"<=": {},
		">=": {},
	}
}

func booleanOperators() map[string]struct{} {
	return map[string]struct{}{
		"&&": {},
		"||": {},
	}
}

func unaryOperators() map[string]struct{} {
	// there are only two lol
	return map[string]struct{}{
		"!": {},
		"-": {},
	}
}

func conditionalKeywords() map[string]struct{} {
	// may later contain switch/match
	return map[string]struct{}{
		"if":      {},
		"else if": {},
		"else":    {},
	}
}

func typeKeywords() map[string]struct{} {
	// taer include map and set
	return map[string]struct{}{
		"int":      {},
		"float":    {},
		"bool":     {},
		"byte":     {},
		"string":   {},
		"function": {},
		"arr":      {},
		"vec":      {},
	}
}

func assignmentKeywords() map[string]struct{} {
	// where keywords are words used to begin statements
	return map[string]struct{}{
		"let": {},
		"mut": {},
	}
}

func iterationKeywords() map[string]struct{} {
	return map[string]struct{}{
		"loop": {},
	}
}

func numbers() map[string]struct{} {
	return map[string]struct{}{
		"0": {},
		"1": {},
		"2": {},
		"3": {},
		"4": {},
		"5": {},
		"6": {},
		"7": {},
		"8": {},
		"9": {},
	}
}

func allKeywords() map[string]struct{} {
	keywords := map[string]struct{}{}
	for k, v := range conditionalKeywords() {
		keywords[k] = v
	}
	for k, v := range typeKeywords() {
		keywords[k] = v
	}
	for k, v := range assignmentKeywords() {
		keywords[k] = v
	}
	for k, v := range iterationKeywords() {
		keywords[k] = v
	}
	return keywords
}
