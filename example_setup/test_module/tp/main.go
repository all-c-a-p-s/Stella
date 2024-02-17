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


func multiply(m tuple4[int, int, int, int], v tuple2[int, int]) tuple2[int, int] {
  return  tuple2[int, int]{v0: m.v0 * v.v0 + m.v1 * v.v1, v1: m.v2 * v.v0 + m.v3 * v.v1}

}
func highest_fibonacci(n int) int {
  var a int  = 0
  var b int  = 1
  for b <= n {
    var temp int  = b
    b = a + b
    a = temp

  }
  return a

}
func zackendorf_representation(n int) bool {
  var ok bool  = true
  if n < 0 {
    ok = false

  } else  {
    var remaining int  = n
    for remaining > 0 {
      var fib int  = highest_fibonacci(remaining)
      fmt.Print( fib )
      fmt.Print( " " )
      remaining = remaining - fib

    }

  }
  return ok

}
func main() {
  var vec  tuple2[int, int] =  tuple2[int, int]{v0: 5, v1: 7}
  var matrix  tuple4[int, int, int, int] =  tuple4[int, int, int, int]{v0: - 1, v1: 0, v2: 0, v3: - 1}
  var result  tuple2[int, int] = multiply(matrix, vec)
  fmt.Print( "result: (" )
  fmt.Print( result.v0 )
  fmt.Print( ", " )
  fmt.Print( result.v1 )
  fmt.Println( ")" + "\n" )
  fmt.Print( "zackendorf representation is: " )
  var ok bool  = zackendorf_representation(5234)
  if ! ok {
    panic( "negative input into zackendorf_representation" )
    fmt.Println( "it is impossible for this code to run" )

  }

}


