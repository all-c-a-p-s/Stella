package main

import (
	"fmt"

	"github.com/all-c-a-p-s/stella/transpiler"
)

func main() {
	var path string = "../test_src/arrays_test.ste"
	fmt.Scanln(&path)
	transpiled := transpiler.TranspileTarget(path)
	formatted := format(transpiled, 2)
	fmt.Println(formatted)
}
