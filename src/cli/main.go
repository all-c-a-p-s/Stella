package main

import (
	"fmt"
	"os"

	"github.com/all-c-a-p-s/stella/transpiler"
)

func main() {
	var path string = "../test_src/syntax_example.ste"
	// fmt.Scanln(&path)
	transpiled := transpiler.TranspileTarget(path)
	formatted := format(transpiled, 2)
	f, err := os.Create("../../test/main.go")
	if err != nil {
		panic("error creating file")
	}

	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			panic("error closing source code file")
		}
	}(f)

	fileSize, err := f.WriteString(formatted)
	if err != nil {
		panic("error writing to file")
	}
	fmt.Println("Transpiled successfully")
	fmt.Printf("Created file ../../test/main.go (%d bytes)\n", fileSize)
}
