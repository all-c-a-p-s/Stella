package transpiler

import (
	"fmt"
	"testing"
)

func TestParseExpression(t *testing.T) {
	// most major errors will get caught here
	// also fairly easy to test

	testScope := Scope{
		arrays:    make(map[string]Array),
		vars:      make(map[string]Variable),
		functions: make(map[string]Function),
		tuples:    make(map[string]Tuple),
	}

	T := ArrayType{
		baseType:   Int,
		dimensions: []int{3},
	}
	foo := Array{
		dataType:   T,
		identifier: "foo",
	}

	boo := Variable{
		dataType:   Int,
		identifier: "boo",
	}

	testScope.vars["boo"] = boo

	testScope.arrays["foo"] = foo

	testScope.functions["bar"] = Function{
		returnType:  Bool,
		identifier:  "bar",
		parameters:  []Variable{},
		arrays:      []Array{foo},
		paramsOrder: []parameterType{ArrayParameter},
	}

	testScope.functions["baz"] = Function{
		returnType:  Bool,
		identifier:  "bar",
		parameters:  []Variable{boo},
		arrays:      []Array{foo},
		paramsOrder: []parameterType{VariableParameter, ArrayParameter},
	}

	var str string = string([]byte{34}) + "string test" + string([]byte{34})
	expr1 := parseExpression(str, 0, &testScope)
	if expr1.dataType != String {
		t.Error("string test failed")
	}

	expr2 := parseExpression("'h'", 0, &testScope)
	if expr2.dataType != Byte {
		t.Error("byte test failed")
	}

	expr3 := parseExpression("(1 + 1) == 2", 0, &testScope)
	if expr3.dataType != Bool {
		t.Error("bool test failed")
	}

	expr4 := parseExpression("foo[2]", 0, &testScope)
	if expr4.dataType != Int {
		t.Error("array indexing test failed")
	}

	expr5 := parseExpression("bar(foo)", 0, &testScope)
	if expr5.dataType != Bool {
		t.Error("array function call test failed")
	}

	expr6 := parseExpression("baz(boo, foo) == (((6.5 + 1.0) > 3.14) || (true == false))", 0, &testScope)
	if expr6.dataType != Bool {
		t.Error("failed ultimate test")
	}

	tStr := string([]byte{34}) + "foo" + string([]byte{34}) + " + " + string([]byte{34}) + "bar" + string([]byte{34})
	expr7 := parseExpression(tStr, 0, &testScope)
	if expr7.dataType != String {
		t.Error("failed string concatenation test")
	}

	_ = parseTupleDeclaration("let tup: (int, float) = (4, 3.14)", 0, &testScope)
	expr8 := parseExpression("tup.1", 0, &testScope)
	if expr8.dataType != Float {
		t.Error("failed tuple indexing test")
	}
}
