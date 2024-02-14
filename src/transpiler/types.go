package transpiler

import (
	"fmt"
)

type (
	primitiveType int
	derivedTypes  int
)

type HasType interface {
	Type() primitiveType
}

const (
	Int primitiveType = iota
	Float
	Bool
	Byte
	String
	IO // used as function return type for main()
)

const (
	arrInt derivedTypes = iota
	arrFloat
	arrBool
	arrByte
	arrString

	vecInt
	vecFloat
	vecBool
	vecByte
	vecString
)

func (p primitiveType) String() string {
	switch p {
	case Int:
		return "int"
	case Float:
		return "float"
	case Bool:
		return "bool"
	case Byte:
		return "byte"
	case String:
		return "string"
	case IO:
		return "IO"
	default:
		panic("Somehow a type not in the enum got passed into primitiveType.String()")
	}
}

func numericType(T primitiveType) bool {
	if T == Int || T == Float {
		return true
	}
	return false
}

func readType(dataType string, lineNum int) primitiveType { // reads data type from assignment
	switch dataType {
	case "int":
		return Int
	case "float":
		return Float
	case "bool":
		return Bool
	case "byte":
		return Byte
	case "string":
		return String
	case "IO":
		return IO
	default:
		panic(fmt.Sprintf("Line %d: data type %s is invalid", lineNum+1, dataType))
	}
}

func getValType(value string, lineNum int) primitiveType {
	if value[0] == '"' && value[len(value)-1] == '"' {
		checkStringVal(value, lineNum)
		return String
	} else if value[0] == 39 && value[len(value)-1] == 39 { // single quotes
		checkByteVal(value, lineNum)
		return Byte
		// integers 0-255 will be parsed as ints and then checked when assigning value
	} else if value == "true" || value == "false" {
		checkBoolVal(value, lineNum)
		return Bool
	} // else must be a number

	foundNum := false

	for i := 0; i < len(value); i++ {
		if _, isNum := numbers()[string(value[i])]; isNum {
			foundNum = true
		}
	}

	if !foundNum {
		panic(fmt.Sprintf("Line %d: unexpected token %s", lineNum+1, value))
	}

	for _, char := range value {
		if char == '.' {
			checkFloatVal(value, lineNum)
			return Float
		}
	}

	checkIntVal(value, lineNum)
	return Int
}

func nextTerm(expression []string, index int, lineNum int) []string { // helper functions for expressionType()
	if index == len(expression)-1 {
		panic(fmt.Sprintf("Line %d: expected another token in expression", lineNum+1))
	}
	bracketCount := 0
	if expression[index+1] != "(" {
		// brackets open a new term of more than one token
		if _, ok := unaryOperators()[expression[index+1]]; ok {
			return nextTerm(expression, index+1, lineNum)
		}
		return []string{expression[index+1]}
		// next token after bracket
	}
	for i := index + 1; i < len(expression); i++ {
		switch expression[i] {
		case "(":
			bracketCount++
		case ")":
			bracketCount--
		}
		if bracketCount == 0 {
			if i == len(expression)-1 {
				return expression[index+1:]
			}
			return expression[index+1 : i+1]
		}
	}
	panic(fmt.Sprintf("Line %d: brackets opened in expression but never closed", lineNum+1))
}

func previousTerm(expression []string, index int, lineNum int) []string { // similar helper function
	if index == 0 {
		panic("previousTerm() was called with no previous term somehow")
	}
	bracketCount := 0
	if expression[index-1] != ")" {
		return []string{expression[index-1]}
	}
	for i := index - 1; i >= 0; i-- {
		switch expression[i] {
		case "(":
			bracketCount++
		case ")":
			bracketCount--
		}
		if bracketCount == 0 {
			return expression[i:index]
		}
	}
	panic(fmt.Sprintf("Line %d: brackets closed in expression but never opened", lineNum+1))
}

func nextOperator(expression []string, index int) (int, error) {
	// returns index of next operator or error if there is none
	bracketCount := 0
	for i := index + 1; i < len(expression); i++ {
		switch expression[i] {
		case "(":
			bracketCount++
		case ")":
			bracketCount--
		}

		if bracketCount == 0 {
			_, ok1 := binaryOperators()[expression[i]]
			_, ok2 := unaryOperators()[expression[i]]
			if ok1 || ok2 {
				return i, nil
			}
		}
	}
	return index, fmt.Errorf("found no next operator in expression")
}

func expressionType(expression []string, lineNum int, currentScope *Scope) primitiveType {
	// NOTE: does not currently support collections
	// collections have a separate function
	// also does not support multi-line expressions
	// which have another separate function

	if len(expression) == 0 {
		panic(fmt.Sprintf("Line %d: Expression is empty", lineNum+1))
	}

	expr := expression // copy made to remove brackets
	if (expr[0] == "(" && expr[len(expr)-1] == ")") || (expr[0] == "{" && expr[len(expr)-1] == "}") {
		var bracketCount int
		var closedBeforeEnd bool

		if len(expr) == 2 {
			panic("Line %d: Expression is empty")
		}
		for i := 0; i < len(expression)-2; i++ { // stop before last index
			switch expression[i] {
			case "(", "{":
				bracketCount++
			case ")", "}":
				bracketCount--
			}
			if bracketCount == 0 {
				closedBeforeEnd = true
				break
			}
		}
		if !closedBeforeEnd { // i.e. last token was the closing bracket
			expr = expr[1 : len(expr)-1]
		}
	}

	if len(expr) == 1 {
		// recursive base case
		_, ok1 := binaryOperators()[expr[0]]
		_, ok2 := unaryOperators()[expr[0]]
		if ok1 || ok2 {
			panic(fmt.Sprintf("Line %d: Expression contains only operators and no values", lineNum+1))
		}

		for i := 0; i < len(expr[0]); i++ {
			if expr[0][i] == '(' {
				fnCall := parseFunctionCall(expr[0], lineNum, currentScope)
				fn := currentScope.functions[fnCall.functionName]
				return fn.returnType
			}
		}

		for i := 0; i < len(expr[0]); i++ {
			// array indexing inside function call would have been caught above
			if expr[0][i] == '[' {
				arrIndex := parseArrayIndexing(expr[0], lineNum, currentScope)
				// might need generics for multi-dimensional arrays
				return arrIndex.dataType.baseType
			}
		}

		if v, ok := (*currentScope).vars[expr[0]]; ok {
			return v.dataType
		}
		return getValType(expr[0], lineNum) // not operator, variable or function
	} else if len(expr) == 2 {
		// also base case as this can only be unary operator and expression of length 1
		switch expr[0] {
		case "-":
			x := parseExpression(expr[1], lineNum, currentScope)
			if x.dataType != Int && x.dataType != Float {
				panic(fmt.Sprintf("Line %d: use of unary operator - with non-numeric data type %v", lineNum+1, x.dataType))
			}
			return x.dataType
		case "!":
			x := parseExpression(expr[1], lineNum, currentScope)
			if x.dataType != Bool {
				panic(fmt.Sprintf("Line %d: use of unary operator - with non-boolean data type %v", lineNum+1, x.dataType))
			}
			return Bool
		default:
			panic(fmt.Sprintf("Line %d: expressions of length 2 tokens must begin with unary operators - or !", lineNum+1))
		}
	}
	typesFound := make(map[primitiveType]struct{}) // Hashset of all types found in expression
	previousOperatorIndex := -1                    // index in expression slice where last term ended.
	// initialise as -1 for first function call

	// match input types of all operators, terms:
	for {
		operatorIndex, err := nextOperator(expr, previousOperatorIndex)
		if err != nil { // no next operator in expression
			break
		}
		previousOperatorIndex = operatorIndex
		operator := expr[operatorIndex]
		switch operator {
		case "-":
			if operatorIndex == 0 {
				next := nextTerm(expr, operatorIndex, lineNum)
				if !numericType(expressionType(next, lineNum, currentScope)) {
					panic(fmt.Sprintf("Line %d: Unary operator '-' found before non numeric type", lineNum+1))
				}
				// typesFound[expressionType(next, lineNum, currentScope)] = struct{}{}
			} else {
				_, prevBinary := binaryOperators()[expr[operatorIndex-1]]
				if prevBinary { // after either numeric operator or comparative operator
					next := nextTerm(expr, operatorIndex, lineNum)
					if !numericType(expressionType(next, lineNum, currentScope)) {
						panic(fmt.Sprintf("Line %d: Unary operator '-' found before non numeric type", lineNum+1))
					}
					// do not add to typesFound if used as unary operator

				} else { // used as binary operator
					previous := previousTerm(expr, operatorIndex, lineNum)
					next := nextTerm(expr, operatorIndex, lineNum)
					if !numericType(expressionType(previous, lineNum, currentScope)) || !numericType(expressionType(next, lineNum, currentScope)) {
						panic(fmt.Sprintf("Line %d: binary opertor '-' used with non-numeric values", lineNum+1))
					}
					typesFound[expressionType(next, lineNum, currentScope)] = struct{}{}
				}
			}
		case "!":
			next := nextTerm(expr, operatorIndex, lineNum)
			if expressionType(next, lineNum, currentScope) != Bool {
				panic(fmt.Sprintf("Line %d: Unary operator '!' used before non-boolean value", lineNum+1))
			}
			// typesFound[Bool] = struct{}{}
		case "+":
			// handle seperately as it can be used with strings as well
			previous := previousTerm(expr, operatorIndex, lineNum)
			next := nextTerm(expr, operatorIndex, lineNum)
			previousType := expressionType(previous, lineNum, currentScope)
			nextType := expressionType(next, lineNum, currentScope)
			if previousType == String && nextType == String {
				typesFound[String] = struct{}{}
			} else {
				if !numericType(previousType) {
					panic(fmt.Sprintf("Line %d: binary operator '%s' used after non-numeric type %v", lineNum+1, expr[operatorIndex], previousType))
				}
				if !numericType(nextType) {
					panic(fmt.Sprintf("Line %d: binary operator '%s' used before non-numeric type %v", lineNum+1, expr[operatorIndex], nextType))
				}
				if previousType != nextType {
					panic(fmt.Sprintf("Lind %d: binary operator '%s' used with both integer and float values", lineNum+1, expr[operatorIndex]))
				}
				typesFound[nextType] = struct{}{}
			}
		case "*", "/":
			previous := previousTerm(expr, operatorIndex, lineNum)
			next := nextTerm(expr, operatorIndex, lineNum)
			previousType := expressionType(previous, lineNum, currentScope)
			nextType := expressionType(next, lineNum, currentScope)
			if !numericType(previousType) {
				panic(fmt.Sprintf("Line %d: binary operator '%s' used after non-numeric type %v", lineNum+1, expr[operatorIndex], previousType))
			}
			if !numericType(nextType) {
				panic(fmt.Sprintf("Line %d: binary operator '%s' used before non-numeric type %v", lineNum+1, expr[operatorIndex], nextType))
			}
			if previousType != nextType {
				panic(fmt.Sprintf("Lind %d: binary operator '%s' used with both integer and float values", lineNum+1, expr[operatorIndex]))
			}
			typesFound[nextType] = struct{}{}
		case "||", "&&":
			previous := previousTerm(expr, operatorIndex, lineNum)
			next := nextTerm(expr, operatorIndex, lineNum)
			previousType := expressionType(previous, lineNum, currentScope)
			nextType := expressionType(next, lineNum, currentScope)
			if previousType != Bool {
				panic(fmt.Sprintf("Line %d: binary operator '%s' used after non-boolean type %v", lineNum+1, expr[operatorIndex], previousType))
			}
			if nextType != Bool {
				panic(fmt.Sprintf("Line %d: binary operator '%s' used before non-boolean type %v", lineNum+1, expr[operatorIndex], nextType))
			}
			// now they must already be the same
			typesFound[Bool] = struct{}{}
		case "==", ">", "<", ">=", "<=", "!=":
			previous := previousTerm(expr, operatorIndex, lineNum)
			next := nextTerm(expr, operatorIndex, lineNum)
			previousType := expressionType(previous, lineNum, currentScope)
			nextType := expressionType(next, lineNum, currentScope)
			if previousType != nextType {
				panic(fmt.Sprintf("Line %d: Binary operator '%s' used with two different types %v and %v", lineNum+1, expr[operatorIndex], previousType, nextType))
			}
			typesFound[Bool] = struct{}{}
		}
	}
	if len(typesFound) == 0 { // shouldn't even be possible to get this lol
		panic(fmt.Sprintf("Line %d: expression has no data type", lineNum+1))
	}
	if len(typesFound) != 1 {
		panic(fmt.Sprintf("Line %d: expression contains more than one data type", lineNum+1))
	}

	var exprType primitiveType
	for k := range typesFound { // all types in expression are the same, and the expression must have a type so this works
		exprType = k
	}
	return exprType
}

func checkIntVal(value string, lineNum int) { // checks to see if int value contains illegal characters/leading zeros etc.
	switch value[0] {
	case '0':
		if len(value) != 1 {
			panic(fmt.Sprintf("Line %d: Integers values cannot have leading zeros", lineNum+1))
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		panic(fmt.Sprintf("Line %d: character %c cannot be part of an integer value", lineNum+1, value[0]))
	}

	for _, char := range value {
		if !(char > 47 && char < 58) { // digits including zero. leading zeros will have been caught above
			panic(fmt.Sprintf("Line %d: character %c cannot be part of an integer value", lineNum+1, char))
		}
	}
}

func checkFloatVal(value string, lineNum int) {
	// check valid float literal
	switch value[0] {
	case '0':
		if !(value[1] == '.') {
			panic(fmt.Sprintf("Line %d: Leading zeros must be followed by decimal point, here it is followed by %c", lineNum+1, value[1]))
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		panic(fmt.Sprintf("Line %d: character %c cannot be part of an float value", lineNum+1, value[0]))
	}

	decimalPointCount := 0
	for _, char := range value {
		if char == '.' {
			decimalPointCount++
		}
		if !(char > 47 && char < 58) { // digits including zero. leading zeros will have been caught above
			if !(char == '.' && decimalPointCount == 1) {
				panic(fmt.Sprintf("Line %d: character %c cannot be part of a float value", lineNum+1, char))
			}
		}
	}
}

func checkBoolVal(value string, lineNum int) {
	// valid bool literal
	if !(value == "true" || value == "false") {
		panic(fmt.Sprintf("Line %d: value '%s' cannot be used as a boolean value", lineNum+1, value))
	}
}

func checkByteVal(value string, lineNum int) {
	// valid byte literal
	if len(value) != 3 {
		panic(fmt.Sprintf("Line %d: single quotes are should be used to enclose single ASCII character, but here there is more than one character inside the quotes", lineNum+1))
	}
	byteVal := []byte(value[1 : len(value)-1])[0]
	if byteVal > 255 {
		panic(fmt.Sprintf("Line %d: value '%s' cannot be used as byte because its ASCII code is over 255", lineNum+1, string(byteVal)))
	}
}

func checkStringVal(value string, lineNum int) {
	// valid string literal
	if !(value[0] == '"' && value[len(value)-1] == '"') {
		panic(fmt.Sprintf("Line %d: '%s' cannot be used as string value", lineNum+1, value))
	}
}
