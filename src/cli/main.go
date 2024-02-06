package main

import (
	"fmt"
	"os"

	"github.com/all-c-a-p-s/stella/transpiler"
)

func main() {
	var path string = "../test_src/prime_factorisation.txt"
	// fmt.Scanln(&path)
	transpiled := transpiler.TranspileTarget(path)
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

	fileSize, err := f.WriteString(transpiled)
	if err != nil {
		panic("error writing to file")
	}
	fmt.Println("Transpiled successfully")
	fmt.Printf("Created file ../../test/main.go (%d bytes)\n", fileSize)
}
