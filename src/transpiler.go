package main

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
	transpiled = transpiled[:len(transpiled)-1] // remove last space
	return transpiled
}

func (D Declaration) transpile() string {
	transpiled := "var "
	transpiled += D.v.identifier + " "
	transpiled += D.v.dataType.String() + " "
	transpiled += " = "
	transpiled += D.e.transpile()
	return transpiled
}

// TODO: finish this to include multi-line expression
func (F Function) transpile() string {
	transpiled := "func "
	transpiled += F.identifier
	transpiled += "("

	for _, p := range F.parameters {
		transpiled += p.identifier
		transpiled += " " + p.dataType.String()
	}

	transpiled += ")"

	transpiled += " " + F.returnType.String()

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
	transpiled += S.condition.transpile()
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
