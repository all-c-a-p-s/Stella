# Stella

The concept of this programming language is that everything is an expression with a type and a value.


variables are declared using the syntax
```
let name: string = "Stella"
```

variables in Stella are immutable by default, and can be made mutable using the ```mut``` keyword

functions are declared like this:
```
function(x: int) -> int = {
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
  a * b --expresion evaluates to this value
}
```

boolean expressions in Stella use brackets for clarity

```
let foo: bool = (1==1)&&(2==2)
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

syntactic characters in Stella each have a single specific purpose

| character | purpose                                                     |
|-----------|-------------------------------------------------------------|
| :         | type annotations                                            |
| {}        | to open and close scopes evaluating to expressions          |
| []        | to contain collections such as arrays, vectors and maps     |
| ()        | to contain parameters of a function                         |
| ,         | separating values                                           |
| <>        | to contain types/sizes of collections                       |
| ->        | is mapped to                                                |


TODO:
-----

collections are types which can be attributed to variables:
```
let nums: arr<int, 5> = [23, 5, 90, 2, 88] --fixed size array
let names: vec<string> = ["bob", "jim", "jeff"] --vector of variable size
let ages: map<string, int> = [
  "bob" -> 12,
  "timmy" -> 50,
  "alex" -> 32,
]

```

the aim is that the simple and consistent syntax should make Stella an approachable and interesting language
