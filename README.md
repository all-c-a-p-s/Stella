# myLang

The concept of this programming language is that everything is an expression with a type and a value.


variables are declared using the syntax
```
let name: type = val
```

variables in myLang are immutable by default, and can be made mutable using the ```mut``` keyword

functions are declared like this:
```
let square: fn(x: int) -> int = {
  x * x
}
```

Expressions in myLang can either be simple single-line expressions such as 
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

selection statements in myLang can be used inside expressions
```
let auth: fn(password: string) -> bool = {
  if password == "secret123" then
    true
  else 
    false
}
```

syntactic characters in myLang each have a single specific purpose

| character | purpose                                                     |
|-----------|-------------------------------------------------------------|
| :         | type annotations                                            |
| {}        | opening and closing scopes evaluating to expressions        |
| []        | containers for collections such as arrays, vectors and maps |
| ()        | containing parameters of a function                         |
| ,         | separating values                                           |
| <>        | containing types/sizes of collections                       |
| ->        | denoting the return type of functions                       |


collections in myLang are types which can be attributed to variables:
```
let nums :: arr<int, 5> = [23, 5, 90, 2, 88] --fixed size array
let names :: vec<string> = ["bob", "jim", "jeff"] --vector of variable size
let ages: map<string, int> = [
  "bob" 12,
  "timmy" 50,
  "alex" 32,
]

```

the aim is that the simple and consistent syntax should make myLang an approachable and interesting language
