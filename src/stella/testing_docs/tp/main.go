package main

import "fmt"
type tuple4[T0 any, T1 any, T2 any, T3 any] struct {
  v0 T0
  v1 T1
  v2 T2
  v3 T3
}

type tuple2[T0 any, T1 any] struct {
  v0 T0
  v1 T1
}


func fibonacci(n int) int {
  var result int  = 1
  if n > 2 {
    result = fibonacci(n-1) + fibonacci(n-2)

  }
  return result

}
func multiply(m tuple4[int, int, int, int], v tuple2[int, int]) tuple2[int, int] {
  return  tuple2[int, int]{v0: m.v0 * v.v0 + m.v1 * v.v1, v1: m.v2 * v.v0 + m.v3 * v.v1}

}
func main() {
  fmt.Println( "works" )

}


