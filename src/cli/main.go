package main

import (
	"fmt"

	"github.com/all-c-a-p-s/stella/transpiler"
)

func main() {
	var path string = "../src.txt"
	// fmt.Scanln(&path)
	transpiled := transpiler.TranspileTarget(path)
	fmt.Println(transpiled)
}
