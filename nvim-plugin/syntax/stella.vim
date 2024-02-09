" Language: Stella
" Maintainer: Sebastiano Rebonato-Scott (https://github.com/all-c-a-p-s)

if exists('b:current_syntax')
    finish
endif

let b:current_syntax = "stella"

"keywords
syntax keyword stellaConditional if else
syntax keyword stellaRepeat      loop
syntax keyword stellaDeclaration let function
syntax keyword stellaStatement   break continue
syntax keyword stellaStorage     mut
syntax keyword stellaType        int float bool string byte
syntax keyword stellaBoolean     true false
syntax keyword stellaTodo        contained TODO FIXME NOTE BUG PERF

"match
syntax match stellaBrackets			 '[\[\]{}=]'
syntax match stellaBrackets			 ' != '
syntax match stellaSign  			 '[+*-]'
syntax match stellaSign 			     ' / '
syntax match stellaSign  			 '->'
syntax match stellaSign  			 '>='
syntax match stellaSign  			 ' > '
syntax match stellaSign  			 '<='
syntax match stellaSign  			 ' < '

syntax match stellaFunctionDeclaration display "\s*\zsfunction\s*\i*\s*<[^>]*>" contains=stellaFunctionName,stellaDeclaration
syntax match stellaFunctionCall        /\w\+\ze\s*(/ contains=stellaDeclaration
syntax match stellaFunctionName        display contained /\s\w\+/

"region
syntax region stellaString              matchgroup=stellaStringDelimiter start=+b"+ skip=+\\\\\|\\"+ end=+"+ 
syntax region stellaString              matchgroup=stellaStringDelimiter start=+"+ skip=+\\\\\|\\"+ end=+"+
syntax region stellaCommentLine         start="//" end="$"   contains=stellaTodo

" Numeric Literals
syntax match       stellaDecimalInt         "\<\d\+\([Ee]\d\+\)\?\>"

syntax match       stellaFloat              "\<\d\+\.\d*\([Ee][-+]\d\+\)\?\>"
syntax match       stellaFloat              "\<\.\d\+\([Ee][-+]\d\+\)\?\>"
syntax match       stellaFloat              "\<\d\+[Ee][-+]\d\+\>"
syntax match       stellaFloat              "[+-]?([0-9]+[.])?[0-9]+"

let s:standardLibraryMacros = ["print", "println"]

for s:standardLibraryMacro in s:standardLibraryMacros
    execute 'syntax match stellaMacros "\v<' . s:standardLibraryMacro . '!"'
endfor

highlight default link stellaString               String
highlight default link stellaBoolean              Boolean
highlight default link stellaSign                 Operator
highlight default link stellaDeclaration          Statement
highlight default link stellaStatement            Statement
highlight default link stellaBrackets             Operator
highlight default link stellaRepeat               Repeat
highlight default link stellaConditional          Conditional
highlight default link stellaOperator             Operator
highlight default link stellaIdentifier           Identifier
highlight default link stellaType                 Type
highlight default link stellaTodo                 TODO
highlight default link stellaStorage              StorageClass
highlight default link stellaCommentLine          Comment
highlight default link stellaMacros               Macros

highlight default link stellaFunctionName        Function
highlight default link stellaFunctionCall        Function

highlight default link stellaDecimalInt    Integer
highlight default link Integer             Number
highlight default link stellaFloat         Float
highlight default link Float               Number
