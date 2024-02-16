# Control Flow in Stella

## Loops

All loops in Stella are declared using the `loop` keyword followed by a boolean condition.

```
let mut i: int = 0
loop i < 10 { //definite iteration
    println!(i)
    i = i + 1
}
```

```
function factorial(x: int) -> int = { //see function documentation for explanation
    let mut n: int = x
    let mut result: int = 1
    loop n > 1 {
        result = result * n
        n = n-1
    }
    result
}
```

```
loop true { //infinite loop, no way this condition's ever gonna be false
    println!("are we nearly there yet?")
}
```

## Selection Statements

Selection statements in Stella execute code is a certain boolean condition is true.

```
function grade_score(mark: int) -> string {
    let mut grade: string = ""
    if mark > 90 {
        grade = "gold"
    } else if mark > 80 {
        grade = "silver"
    } else {
        grade = "bronze"
    }
    grade
}
```
