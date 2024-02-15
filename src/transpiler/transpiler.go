package transpiler

import (
	"fmt"
)

//line which is just "}"
type ScopeCloser struct {
	closer string
}

type Transpileable interface {
	transpile() string
}

func generateTupleCode(n int) string {
	// generates struct that needs to be added to top of Go file
	// when tuple of size n is used
	var transpiled string
	transpiled += "type tuple"
	transpiled += fmt.Sprintf("%d", n) + "["
	for i := 0; i < n; i++ {
		transpiled += "T" + fmt.Sprintf("%d", i) + " any"
		if i != n-1 {
			transpiled += ", "
		}
	}
	transpiled += "]" + " struct {"
	transpiled += "\n"
	for i := 0; i < n; i++ {
		transpiled += "v" + fmt.Sprintf("%d", i) + " T" + fmt.Sprintf("%d", i)
		transpiled += "\n"
	}
	transpiled += "}"

	return transpiled
}

func (E Expression) transpile() string {
	items := E.items
	var transpiled string
	for _, item := range items {
		if isTupleIndexing(item) {
			transpiled += transpileTupleIndexing(item)
		} else {
			transpiled += item
		}
		transpiled += " "
	}
	if len(transpiled) != 0 {
		transpiled = transpiled[:len(transpiled)-1] // remove last space
	}
	return transpiled
}

func (D Declaration) transpile() string {
	transpiled := "var "
	transpiled += D.v.identifier + " "
	if D.v.dataType != Float {
		transpiled += D.v.dataType.String() + " "
	} else {
		transpiled += D.v.dataType.String() + "64" + " "
	}
	transpiled += " = "
	transpiled += D.e.transpile()
	return transpiled
}

func (A ArrayDeclaration) transpile() string {
	transpiled := "var "
	transpiled += A.arr.identifier + " "

	transpiled += "["
	transpiled += fmt.Sprintf("%d", A.arr.dataType.dimensions[0])
	transpiled += "]"
	if A.arr.dataType.baseType != Float {
		transpiled += A.arr.dataType.baseType.String()
	} else {
		transpiled += A.arr.dataType.baseType.String() + "64"
	}
	transpiled += " = "
	transpiled += A.expr.transpile()
	return transpiled
}

func (A ArrayValue[primitiveType]) transpile() string {
	var transpiled string
	transpiled += "[" + fmt.Sprintf("%d", A.length) + "]"

	if A.baseType != Float {
		transpiled += A.baseType.String()
	} else {
		transpiled += A.baseType.String() + "64"
	}
	transpiled += "{"
	for i, elem := range A.elements {
		transpiled += elem.transpile()
		if i != len(A.elements)-1 {
			transpiled += ", "
		}
	}
	transpiled += "}"
	return transpiled
}

func (F Function) transpile() string {
	transpiled := "func "
	transpiled += F.identifier
	transpiled += "("

	var varCount, arrCount, tupCount int

	for i, t := range F.paramsOrder {
		if t == VariableParameter {
			p := F.parameters[varCount]
			transpiled += p.identifier
			if p.dataType != Float {
				transpiled += " " + p.dataType.String()
			} else {
				transpiled += " " + p.dataType.String() + "64"
			}
			if i != len(F.paramsOrder)-1 {
				transpiled += ", "
			}
			varCount++
		} else if t == TupleParameter {
			t := F.tuples[tupCount]
			transpiled += t.identifier
			transpiled += " tuple"
			transpiled += fmt.Sprintf("%d", len(t.pattern.dataTypes))
			transpiled += "["
			for i, T := range t.pattern.dataTypes {
				switch T {
				case Int:
					transpiled += "int"
				case Float:
					transpiled += "float64"
				case Bool:
					transpiled += "bool"
				case Byte:
					transpiled += "byte"
				case String:
					transpiled += "string"
				}
				if i != len(t.pattern.dataTypes)-1 {
					transpiled += ", "
				}
			}
			transpiled += "]"
			if i != len(F.paramsOrder)-1 {
				transpiled += ", "
			}
			tupCount++
		} else {
			arr := F.arrays[arrCount]
			transpiled += arr.identifier
			transpiled += " "
			transpiled += "[" + fmt.Sprintf("%d", arr.dataType.dimensions[0]) + "]"
			transpiled += arr.dataType.baseType.String()
			if arr.dataType.baseType == Float {
				transpiled += "64"
			}
			if i != len(F.paramsOrder)-1 {
				transpiled += ", "
			}
			arrCount++
		}
	}

	transpiled += ")"

	if F.returnDomain == tuple {
		T := F.tupleReturnType
		transpiled += " " + "tuple"
		transpiled += fmt.Sprintf("%d", len(T.dataTypes))
		transpiled += "["
		for i, dataType := range T.dataTypes {
			switch dataType {
			case Int:
				transpiled += "int"
			case Float:
				transpiled += "float64"
			case Bool:
				transpiled += "bool"
			case Byte:
				transpiled += "byte"
			case String:
				transpiled += "string"
			}
			if i != len(T.dataTypes)-1 {
				transpiled += ", "
			}
		}
		transpiled += "]"
	} else if F.returnDomain == derived {
		if len(F.derivedReturnType.dimensions) == 0 {
			panic("shouldn't be possible to panic here 🙏")
		}
		transpiled += " " + "[" + fmt.Sprintf("%d", F.derivedReturnType.dimensions[0]) + "]" + F.derivedReturnType.baseType.String()
		if F.derivedReturnType.baseType == Float {
			transpiled += "64"
		}
	} else {
		if F.returnType != IO {
			transpiled += " " + F.returnType.String()
			if F.returnType == Float {
				transpiled += "64"
			}
		}
	}

	transpiled += " {"

	return transpiled
}

func (L Loop) transpile() string {
	transpiled := "for "
	transpiled += L.condition.transpile()
	transpiled += " {"
	return transpiled
}

func (B BreakStatement) transpile() string {
	switch B.T {
	case Break:
		return "break"
	case Continue:
		return "continue"
	}
	panic("should be literally impossible for transpiler to ever panic here lol")
}

func (S SelectionStatement) transpile() string {
	var transpiled string
	switch S.selectionType {
	case If:
		transpiled = "if "
	case ElseIf:
		transpiled = "} else if "
	case Else:
		transpiled = "} else "
	}
	transpiled += S.condition.transpile() + " "
	transpiled += "{"
	return transpiled
}

func (A Assignment) transpile() string {
	var transpiled string
	transpiled += A.v.identifier
	transpiled += " = "
	transpiled += A.e.transpile()
	return transpiled
}

func (A ArrayAssignment) transpile() string {
	var transpiled string
	transpiled += A.arr.identifier
	transpiled += " = "
	transpiled += A.expr.transpile()
	return transpiled
}

func (A ArrayIndexAssignment) transpile() string {
	var transpiled string
	transpiled += A.arrIndex.arrayID
	transpiled += "["
	transpiled += A.arrIndex.index.transpile()
	transpiled += "]"

	transpiled += " = "
	transpiled += A.value.transpile()
	return transpiled
}

func (A ArrayExpression) transpile() string {
	if len(A.literal.values) > 0 {
		// fine as there are no operators that work on arrays
		return A.literal.transpile()
	}
	return A.stringValue
	// fine as function call can't contain literals
}

func (B BaseArray) transpile() string {
	var transpiled string
	transpiled += "[" + fmt.Sprintf("%d", B.length) + "]"

	if B.dataType != Float {
		transpiled += B.dataType.String()
	} else {
		transpiled += B.dataType.String() + "64"
	}
	transpiled += "{"
	for i, elem := range B.values {
		transpiled += elem.transpile()
		if i != len(B.values)-1 {
			transpiled += ", "
		}
	}
	transpiled += "}"

	return transpiled
}

func (T TupleLiteral) transpile() string {
	// necessary struct and interface already created above
	var transpiled string
	transpiled += " tuple"
	transpiled += fmt.Sprintf("%d", len(T.values))
	transpiled += "["
	for i, expr := range T.values {
		switch expr.dataType {
		case Int:
			transpiled += "int"
		case Float:
			transpiled += "float64"
		case Bool:
			transpiled += "bool"
		case Byte:
			transpiled += "byte"
		case String:
			transpiled += "string"
		}
		if i != len(T.values)-1 {
			transpiled += ", "
		}
	}
	transpiled += "]"
	transpiled += "{"
	for i, e := range T.values {
		transpiled += "v" + fmt.Sprintf("%d", i) + ":" + " "
		transpiled += e.transpile()
		if i != len(T.values)-1 {
			transpiled += ", "
		}
	}
	transpiled += "}"
	return transpiled
}

func (F FunctionCall) transpile() string {
	var transpiled string
	transpiled += F.functionName + "("
	var varCount, arrCount, tupCount int
	for i := 0; i < len(F.order); i++ {
		switch F.order[i] {
		case VariableParameter:
			transpiled += F.parameters[varCount].transpile()
			varCount++
		case ArrayParameter:
			transpiled += F.arrays[arrCount].identifier
			arrCount++
		case TupleParameter:
			transpiled += F.tuples[tupCount].identifier
			tupCount++
		}
		if i != len(F.order)-1 {
			transpiled += ", "
		}
	}
	transpiled += ")"
	return transpiled
}

func (T TupleExpression) transpile() string {
	if T.exprType == Literal {
		return T.literal.transpile()
	} else if T.exprType == FnCall {
		return T.fnCall.transpile()
	}
	return T.t.identifier
}

func (T TupleDeclaration) transpile() string {
	transpiled := "var "
	transpiled += T.t.identifier + " "
	transpiled += " tuple"
	transpiled += fmt.Sprintf("%d", len(T.t.pattern.dataTypes))
	transpiled += "["
	for i, dataType := range T.t.pattern.dataTypes {
		switch dataType {
		case Int:
			transpiled += "int"
		case Float:
			transpiled += "float64"
		case Bool:
			transpiled += "bool"
		case Byte:
			transpiled += "byte"
		case String:
			transpiled += "string"
		}
		if i != len(T.t.pattern.dataTypes)-1 {
			transpiled += ", "
		}
	}
	transpiled += "]"
	transpiled += " = "
	transpiled += T.e.transpile()
	return transpiled
}

func (T TupleAssignment) transpile() string {
	transpiled := T.t.identifier
	transpiled += " = "
	transpiled += T.e.transpile()
	return transpiled
}

func (T TupleIndexing) transpile() string {
	transpiled := T.t.identifier
	transpiled += "."
	transpiled += "v" + fmt.Sprintf("%d", T.i)
	return transpiled
}

func (M Macro) transpile() string {
	var transpiled string
	switch M.T {
	case Print:
		transpiled += "fmt.Print"
	case Println:
		transpiled += "fmt.Println"
	case Panic:
		transpiled += "panic"
	default:
		panic("macro not supported by transpile()")
	}

	transpiled += M.value.transpile()
	return transpiled
}

func (s Scope) transpile() string {
	var transpiled string

	for _, item := range s.items {
		T := typeOfItem(item)

		if T == "Expression" || T == "ArrayExpression" || T == "TupleExpression" {
			transpiled += "return " + item.transpile()
			// only way for an expression to come alone
		} else {
			transpiled += item.transpile()
		}
		transpiled += "\n"
	}

	return transpiled
}

func (S ScopeCloser) transpile() string {
	return S.closer
}
