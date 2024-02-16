# Stella

## Name

- Acronym: Strongly Typed Expression Lightweight LAnguage
- stella means star in Italian
- named after flower found in the alps "Stella Alpina"

## Purpose

Stella has 3 main aims:

- to be simple and accessible to learn
- to include features making it easy to write bug-free code
- to be fast, making it practical to be used

The combination of these make it's best use-case scientific modelling. It is easy to learn and it's code executes very quickly. (Note: the project is not developed enough to be a realistic option for professionals to use).

## Features

- transpiles into Go, which is a very fast compiled language
- pointers/references are disallowed
- pass by value only
- These two make Stella impossible to use for embedded/systems programming. However, for the intended use case of data modelling, these encourage a simple flow of data through functions as the program executes.
- C-like syntax. This will make it familiar to people with experience in languages such as Python, Javascript, and Go
- static typing. This makes it easier for the developer to create a complex and bug-free program, while also enabling the transpiling to the fast compiled language Go.
- mathematically intuitive syntax (subjective). Some of Stella's syntax is inspired by functional languages such as Haskell (e.g. -> for function return type).

## Support:

![demo](demo.PNG)

I made some syntax highlighting plugins. They are not perfect, but should be useful for people who want to try out the language.

**NeoVim**: [Stella-nvim](https://github.com/all-c-a-p-s/Stella-nvim)

**VSCode**: [Stella-Lang](https://marketplace.visualstudio.com/items?itemName=StellaLang.stella-lang)

## Contributing

Not currently accepting contributions as this is a school project.
