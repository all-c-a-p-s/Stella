" Language: Stella
" Maintainer: Sebastiano Rebonato-Scott (https://github.com/all-c-a-p-s)

if exists('b:current_syntax')
    finish
endif

let b:current_syntax = "stella"

"keywords
syn keyword stellaConditional if else
syn keyword stellaRepeat      loop
syn keyword stellaKeyword     let function
syn keyword stellaStatement   break continue
syn keyword stellaStorage     mut
syn keyword stellaType        int float bool string byte
syn keyword stellaBoolean     true false
syn keyword stellaTodo        contained TODO FIXME NOTE BUG PERF

"match
syn match stellaBrackets			 '[\[\]{}=]'
syn match stellaBrackets			 ' != '
syn match stellaSign  			 '[+*-]'
syn match stellaSign 			     ' / '
syn match stellaSign  			 '->'
syn match stellaSign  			 '>='
syn match stellaSign  			 ' > '
syn match stellaSign  			 '<='
syn match stellaSign  			 ' < '

"region
syn region stellaString              matchgroup=stellaStringDelimiter start=+b"+ skip=+\\\\\|\\"+ end=+"+ 
syn region stellaString              matchgroup=stellaStringDelimiter start=+"+ skip=+\\\\\|\\"+ end=+"+
syn region stellaCommentLine         start="//" end="$"   contains=stellaTodo

" Integers
syn match       stellaDecimalInt         "\<\d\+\([Ee]\d\+\)\?\>"

syn match       stellaFloat              "\<\d\+\.\d*\([Ee][-+]\d\+\)\?\>"
syn match       stellaFloat              "\<\.\d\+\([Ee][-+]\d\+\)\?\>"
syn match       stellaFloat              "\<\d\+[Ee][-+]\d\+\>"
syn match       stellaFloat              "[+-]?([0-9]+[.])?[0-9]+"

highlight def link stellaString               String
highlight def link stellaBoolean              Boolean
highlight def link stellaSign                 Type
highlight def link stellaKeyword              Keyword
highlight def link stellaStatement            Statement
highlight def link stellaBrackets             Type
highlight def link stellaRepeat               Repeat
highlight def link stellaConditional          Conditional
highlight def link stellaOperator             Operator
highlight def link stellaIdentifier           Identifier
highlight def link stellaType                 Type
highlight def link stellaTodo                 TODO
highlight def link stellaStorage              StorageClass
highlight def link stellaCommentLine          Comment

highlight def link stellaDecimalInt    Integer
highlight def link Integer             Number
highlight def link stellaFloat         Float
highlight def link Float               Number
