package main

import "fmt"

func fibonacci(n int) int {
	var to_return int = 1
	if n > 2 {
		to_return = fibonacci(n-1) + fibonacci(n-2)
	}
	return to_return
}

func factorial(x int) int {
	var to_return int = 1
	if x > 1 {
		to_return = x * factorial(x-1)
	}
	return to_return
}

func exp2(x int) int {
	var to_return int = 1
	if x > 0 {
		to_return = 2 * exp2(x-1)
	}
	return to_return
}

func celsius_to_fahrenheit(temp float64) float64 {
	return temp*1.8 + 32.0
}

func main() {
	fmt.Print("hello world")
}
