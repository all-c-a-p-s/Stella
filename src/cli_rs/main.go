package main

import "fmt"

func max(arr [10]float64) float64 {
	var answer float64 = 0.0
	var i int = 1
	for i < 10 {
		if arr[i] < 0.0 {
			fmt.Println("negative input")
			break

		} else if arr[i] > answer {
			answer = arr[i]
		}
		i = i + 1

	}
	return answer
}

func fib(n int) int {
	var x int = 1
	if n < 2 {
		x = 1
	} else {
		x = fib(x-1) + fib(x-2)
	}
	return x
}

func main() {
	var arr [10]float64 = [10]float64{0.2, 2.718, 3.14, 1.618, -5.0, 16.0, 44.4, 23.4, 0.01, -11.0}
	fmt.Println("Hello from Stella")
	fmt.Print("Tenth fibonacci number is: ")
	fmt.Println(fib(10))
	fmt.Print("Max of array is: ")
	fmt.Println(max(arr))
}
