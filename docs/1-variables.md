# Variables In Stella

There are three types of variable in Stella: primitive variables, derived variables and product variables.

## Primitives

```
let name: string = "Stella" //type annotation after the colon is mandatory
let mut version: int = 1 //all variables must be assigned a value
//  ^ the 'mut' keyword indicates that the variables is mutable
// variables without the 'mut' keyword are immutable constants
```

Primitive variables in Stella are values of one of the following primitive data types:
| Type | Meaning |
|--------|---------------------------------------------------------------------------------------------|
| int | an integer number (either 32 bit or 64 bit) depending on the user's system. e.g. 12, -3, 0 |
| float | a 64-bit floating point (decimal) number e.g. 3.0, -55.5 |
| bool | a boolean expression (true/false) e.g. true, false, (5 > 3) |
| byte | an ASCII-encoded character e.g. 'A', 'a', ' ' |
| string | a string of Unicode characters e.g. "Hello from Stella ✨!" |

## Derived

Derived data types are defined in terms of primitive types.

### Arrays

Arrays are Stella's core implementation of derived data types. They are fixed-size collections of values which are all of the same data type.

```
let nums: float[3] = [3.14, 2.71, 1.62]
//        ^ type annotation with size in brackets
// [3.14, 2.71, 1.62] is an array literal
```

## Product

A boolean variable can have 2 possible values (true or false)
A tuple of two booleans can have 2\*2 = 4 possible values (FF, FT, TF, TT). Hence, data structures with multiple fields (of fixed types) are called product types.

### Tuples

```
let person1: (string, int) = ("Jim", 21)
let identity_matrix(int, int, int, int) = (1, 0, 0, 1)
print!(person.0 + "'s age is ") //tuples are indexed using the .n syntax (zero-indexed)
println!(person1.1)
```
