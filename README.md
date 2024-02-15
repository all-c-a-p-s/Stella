# Stella

### Stella (Strongly Typed Expressive Lightweight LAnguage) has 3 main aims:

- to be simple to learn and accessible, and easy to write code in
- to include features making it easy to write bug-free code
- to be fast, with similar performance to the garbage-collected Go (which it currently transpiles to)

![demo](demo.PNG)

## Support:

I made some syntax highlighting plugins. They are not perfect, but should be useful for people who want to try out the language.

**NeoVim**: [Stella-nvim](https://github.com/all-c-a-p-s/Stella-nvim)

**VSCode**: [Stella-Lang](https://marketplace.visualstudio.com/items?itemName=StellaLang.stella-lang)

## Documentation:

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
  a * b // expression evaluates to this value
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

Loops in Stella use the `loop` keyword with a boolean condition like this

```
let mut i: int = 0
loop i < 10 {
  println!(i)
  i = i + 1
}
```

Arrays in Stella are collections of primitive data types. They can be created like this:

```
let nums: int[5] = [1, 2, 3, 4, 5]
```

and indexed like this

```
let first: int = nums[0]
```


Tuples in Stella are collections of fixed data types (which can be different) in a fixed order:

```
let person1: (string, int) = ("mark", 60)
print!(person1.0 + "'s age is ") //tuples are indexed using the .n syntax (zero-indexed)
println!(person1.1)
```

## TODO:

vectors are arrays of dynamic size (heap allocated):

```
let names: string[] = ["tim", "sarah", "sam"] // vector of variable size
names = append(names, "sid")
```

the aim is that the simple and consistent syntax should make Stella an approachable and interesting language
