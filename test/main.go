package main

import "fmt"

func divisible(n int, divisor int, previous int) bool {
var to_return bool  = false
if ( previous + divisor ) > n {
to_return = false

} else if ( previous + divisor ) == n {
to_return = true

} else  {
to_return = divisible(n, divisor, previous + divisor)

}
return to_return

}
func next_divisor(n int, divisor int) int {
var to_return int  = 0
if divisible(n, divisor, 0) {
to_return = divisor

} else  {
to_return = next_divisor(n, divisor+1)

}
return to_return

}
func divide(n int) int {
var to_return int  = 0
if n == 1 {
fmt.Println( n )

} else  {
var d int  = next_divisor(n, 2)
fmt.Println( d )
var divided int  = n / d
var res int  = divide(divided)
to_return = res

}
return to_return

}
func main() {
var num_to_factorise int  = 228944
var res int  = divide(num_to_factorise)
if res == - 1 {
fmt.Print( res )

}

}
