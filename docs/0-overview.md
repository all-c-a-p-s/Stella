# Stella - Overview

Stella has 3 main aims:

- to be simple and accessible to learn
- to include features making it easy to write bug-free code
- to be fast, making it practical to be used

## Features

- transpiles into Go, which is a very fast compiled language
- pointers/references are disallowed
- pass by value only
- These two make Stella impossible to use for embedded/systems programming. However, for the intended use case of data modelling, these encourage a simple flow of data through functions as the program executes.
- C-like syntax. This will make it familiar to people with experience in languages such as Python, Javascript, and Go
- static typing. This makes it easier for the developer to create a complex and bug-free program, while also enabling the transpiling to the fast compiled language Go.
- mathematically intuitive syntax (subjective). Some of Stella's syntax is inspired by functional languages such as Haskell (e.g. -> for function return type).
