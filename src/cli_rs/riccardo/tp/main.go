package main

import "fmt"

func square(x int) int {
  return x * x

}
func factorial(x int) int {
  var toreturn int  = 1
  if ( x != 0 ) && ( x != 1 ) {
    var j int  = 1
    for j <= x {
      toreturn = toreturn * j
      j = j + 1

    }

  }
  return toreturn

}
func main() {
  fmt.Println( "seba is clever" )
  fmt.Print( "The square of the number you thought is " )
  fmt.Println( square(5) )
  var i int  = 2
  for i < 10 {
    fmt.Print( "if i is " )
    fmt.Print( i )
    fmt.Print( " then the factorial of i is " )
    fmt.Println( factorial(i) )
    i = i + 1

  }

}


