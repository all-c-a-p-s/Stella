package main

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

	return transpiled
}
