package main

import (
	"fmt"
	"os"

	"github.com/all-c-a-p-s/stella/transpiler"
)

func main() {
	ok, err := os.ReadFile("metadata.txt")
	if err != nil {
		panic(err)
	}

	path := string(ok)

	transpiled := transpiler.TranspileTarget(path)
	formatted := format(transpiled, 2)
	fmt.Println(formatted)
}
