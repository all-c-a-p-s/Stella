# Functions In Stella

Functions in Stella are defined using the syntax

```typescript
function fibonacci(n: int) -> int = { //parameters and return type must have explicit type annotation
    let mut result: int = 1
    if n > 2 {
        result = fibonacci(n-1) + fibonacci(n-2) //recursive calls
    }
    result // return value is the expression on the last line
}
```

The idea of this syntax is that it reads like an English sentence: The function fibonacci maps (->) n, an integer, to an integer equal to {...}. All functions is stella must open a block (a multi-line expression) with curly brackets.

This is the reason for the mandatory return value as the expression on the last line of the function.

Functions can take primitive, derived or product variables as parameters (and as return type).

```typescript
function multiply(m: (int, int, int, int), v: (int, int)) -> (int, int) = {
    (m.0 * v.0 + m.1 * v.1, m.2 * v.0 + m.3 * v.1)
}
```

All stella files must contain a special main() function which is the only function which is allowed to have the return type IO (input/output).

```typescript
function main() -> IO = {
    println!("this is the main function")
    println!("code must be called from here in order to be executed")
}
```
