package transpiler

import "fmt"

//line which is just "}"
type ScopeCloser struct {
	closer string
}

type Transpileable interface {
	transpile() string
}

func (E Expression) transpile() string {
	items := E.items
	var transpiled string
	for _, item := range items {
		transpiled += item
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

	for _, p := range F.parameters {
		transpiled += p.identifier
		if p.dataType != Float {
			transpiled += " " + p.dataType.String()
		} else {
			transpiled += " " + p.dataType.String() + "64"
		}
	}

	transpiled += ")"

	transpiled += " " + F.returnType.String()
	if F.returnType == Float {
		transpiled += "64"
	}

	transpiled += " {"

	return transpiled
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

func (s Scope) transpile() string {
	var transpiled string

	for _, item := range s.items {
		T := typeOfItem(item)

		if T == "Expression" {
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
