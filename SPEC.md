# Specification for Rouph
 This is a reference manual for the Rouph programming language.
 The grammar is compact and simple to parse.

## Notation

The syntax is specified using Extended Backus-Naur Form (EBNF):

```
program    = function*
function   = statement*
statement  = expression ";"
            | "{" statement* "}"
expresion  = assign
assign     = equality("=" assign)?
equality   = relational ("==" relational | "!=" relational)*
relational = add("<" add | "<=" add | ">" add | ">=" add)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary| "%" unary )*
unary      = ("+" | "-" | "*" | "&" )? primary
primary    = num | identifier | "(" expr ")"
```
## Literals
Rouph includes 4 literals:
- int
    -  the set of all unsigned 31-bit integers (0 to 2147483647)
- char
    - A integer value identifying an ASCII code point. 
- Array
    - A descriptor for a contiguous segment of an underlying array and provides access to a numbered sequence of elements from that array. 
- Pointer
    - A pointer type denotes the set of all pointers to variables of a given type, called the base type of the pointer. 
## Variables
A variable is a storage location for holding a value. The set of permissible values is determined by the variable's type. 

```
let x:int;
let y:int = 1;
let z:char* = "abcd";
```

## Function types
A function type denotes the set of all functions with the same parameter and result types. 

```
func()
func(x: int):int
func(y: int, z:int):int
```

## if statements
"If" statements specify the conditional execution of two branches according to the value of a boolean expression. If the expression evaluates to true, the "if" branch is executed, otherwise, if present, the "else" branch is executed. 

```
if x == 1 {
    return 0;
}
```
## for statements
A "for" statement specifies repeated execution of a block.

```
let i:int;
let j:int=1;
for i=1; i <21; i= i+1 {
    i = i*j;
}
```

```
let i:int;
for i=1;;i= i+1 {
    if i == 10 {
        return i;
    }
}
```
## Built-in functions
Built-in functions are predeclared. 
### put
Writes to th standard output, accepts character, string, integer literals.
```
put("Hello");
```
```
let i:int=0;
put(i);
```
```
let s:char* = "abcd";
put(s);
```