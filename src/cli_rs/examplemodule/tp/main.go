package main

import "fmt"
type tuple2[T0 any, T1 any] struct {
  v0 T0
  v1 T1
}

type tuple1[T0 any] struct {
  v0 T0
}


func square(x int) int {
  return x * x

}
func invert_and_multiply(x tuple2[int, int], y int) tuple2[int, int] {
  return  tuple2[int, int]{v0: x.v1 * y, v1: x.v0 * y}

}
func add_one_to_all(nums [10]int) [10]int {
  var res [10]int = nums
  var i int  = 0
  for i < 10 {
    res[i] = res[i] + 1

  }
  var foo  tuple1[int] =  tuple1[int]{v0: 5}
  res[9] = foo.v0
  return   res

}
func main() {
  fmt.Println( "hi" )
  var foo  tuple1[int] =  tuple1[int]{v0: 2}
  var nums [10]int = [10]int{0, 1, 2, 3, 4, foo.v0, 6, 7, 8, 9}
  var added [10]int = add_one_to_all([0, 1, 2, 3, 4, foo.0, 6, 7, 8, 9])
  fmt.Println( added[0] )

}


