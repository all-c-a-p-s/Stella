# Stella

### Stella (Strongly Typed Expressive Lightweight LAnguage) has 3 main aims:

- to be simple to learn and accessible, and easy to write code in
- to include features making it easy to write bug-free code
- to be fast, with similar performance to the garbage-collected Go (which it currently transpiles to)

variables are declared using the syntax

```
let name: string = "Stella"
```

variables in Stella are immutable by default, and can be made mutable using the `mut` keyword

functions are declared like this:

```
function square(x: int) -> int = {
  x * x
}
```

Expressions in Stella can either be simple single-line expressions such as

```
x * x
```

or multi-line expressions evaluating to a single value:

```
{
  let a: int = 5
  let b: int = 10
  a * b --expression evaluates to this value
}
```

boolean expressions in Stella use brackets for clarity

```
let auth: bool = (username == "marvin") && (password == "secret123")
```

selection statements in Stella are used with the if keyword and a boolean expression

```
let mut ok: bool = false
if password == "secret123" {
  ok = true
} else if password == "let me in please" {
  ok = true
} else {
  ok = false
}
```

Single-dimension arrays in Stella are currently supported in Stella, with support for multi-dimensional arrays coming soon.

Arrays can be created like this:
```
let nums: int[5] = [1, 2, 3, 4, 5]
```

## TODO:

vectors are arrays of dynamic size (heap allocated):

```
let names: string[] = ["tim", "sarah", "sam"] --vector of variable size
```

the aim is that the simple and consistent syntax should make Stella an approachable and interesting language
