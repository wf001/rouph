package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var standard_lib = []string{"put"}
var reserve_sig = []string{"return", "if", "else", "for", "func", "let", "int", "sizeof", "char"}
var eq_rel_op = []string{"==", "!=", ">=", "=<"}
var op = "+-*/()><;={},&[]%:"

type TokKind int

const (
	TK_KIND_RESERVED TokKind = iota + 1
	TK_IDENT
	TK_STR
	TK_KIND_NUM
	TK_KIND_EOF
	TK_STDLIB
)

type Token struct {
	Kind     TokKind
	Next     *Token
	Val      string
	Str      string
	Contents string
	ContLen  int
}

func printToken(tok *Token) {
	if DEBUG {
		fmt.Println("============= print token =============")
		for ; tok != nil; tok = tok.Next {
			Info("##\x1b[36m tok %p\x1b[0m\n", tok)
			Info("%+v\n", tok)
		}
		fmt.Println("==========================")
	}
}
func Info(s string, v interface{}) {
	if DEBUG {
		_, file, line, _ := runtime.Caller(1)
		reg := "[/]"
		files := regexp.MustCompile(reg).Split(file, -1)
		fmt.Printf("%s %d| ", files[len(files)-1], line)
		fmt.Printf(s, v)
	}
}
func sep() {
	fmt.Println("--------------------------------")
}

func newToken(kind TokKind, cur *Token, val string) *Token {
	tok := new(Token)
	tok.Kind = kind
	if kind == TK_IDENT {
		tok.Str = val
	} else {
		tok.Val = val
	}
	cur.Next = tok
	return tok
}
func TokenizeHandler(arg string) *Token {
	var arg_arr []string
	arg_arr = Scan(arg, 0, len(arg), arg_arr, "", "")

	head := new(Token)
	head.Next = nil
	head.Kind = -1
	cur := head

	//tokenize
	Info("arg_arr:%#v\n", arg_arr)
	for _, s := range arg_arr {
		// reserved signature
		if _, res := isReserved(s); res {
			cur = newToken(TK_KIND_RESERVED, cur, s)
			// return statement
		} else if s == "return" {
			cur = newToken(TK_KIND_RESERVED, cur, s)
			// standard library
		} else if _, res := startWith(s, standard_lib); res {
			cur = newToken(TK_STDLIB, cur, s)
			// identifier
		} else if isAlpha(rune(s[0])) {
			cur = newToken(TK_IDENT, cur, s)
			// string literal
		} else if string(s[0]) == "\"" {
			cur = newToken(TK_STR, cur, s)
			cur.ContLen = len(string(s)) - 1
			cur.Contents = string(s)[1:cur.ContLen] + string(rune(0))
		} else {
			cur = newToken(TK_KIND_NUM, cur, s)
		}
	}
	cur = newToken(TK_KIND_EOF, cur, "")

	return head
}

func Scan(input string, i int, n int, arr []string, number string, valName string) []string {
	if i == n {
		arr = appendIfExists(arr, number, valName)
		return arr
	}
	Info("input %s\n", string(input[i]))
	// if input[i] is space, it skipped
	if string(input[i]) == " " {
		arr = appendIfExists(arr, number, valName)
		number, valName = "", ""
		i += 1
		// if input[i] is '"'
	} else if string(input[i]) == "\"" {
		arr = appendIfExists(arr, number, valName)
		number, valName = "", ""
		strVal := string(input[i])
		for {
			i += 1
			if string(input[i]) == "\"" {
				break
			}
			if string(input[i]) == "\\" {
				i += 1
				strVal += string(getEscapeChar(string(input[i])))
			} else {
				strVal += string(input[i])
			}
		}
		strVal += string(input[i])
		arr = append(arr, strVal)
		i += 1
		// if input[i] is reserved signature
	} else if sig, res := startWith(string(input[i:]), reserve_sig); res &&
		len(valName) == 0 &&
		!isAlphaNum(rune(input[i+len(sig)])) {
		arr = appendIfExists(arr, number, valName)
		arr, number, valName = append(arr, sig), "", ""
		i += len(sig)
		// if input[i] is standard library
	} else if sig, res := startWith(string(input[i:]), standard_lib); res &&
		len(valName) == 0 {
		arr = appendIfExists(arr, number, valName)
		arr, number, valName = append(arr, sig), "", ""
		i += len(sig)
		// if input[i] is Equality or Relational operator
	} else if op, res := startWith(string(input[i:]), eq_rel_op); res {
		arr = appendIfExists(arr, number, valName)
		arr, number, valName = append(arr, op), "", ""
		i += 2
		// if input[i] is single-letter operator
	} else if op, res := isReserved(string(input[i])); res {
		arr = appendIfExists(arr, number, valName)
		arr, number, valName = append(arr, op), "", ""
		i += 1
		// if input[i] is alphabet or number, it taken as a part of variable name.
		// note that ahead of val must be alphabet.
	} else if isAlpha(rune(input[i])) ||
		(isAlphaNum(rune(input[i])) && valName != "") {
		number = ""
		valName += string(input[i])
		i += 1
		// otherwise, input[i] must be number
	} else {
		number += string(input[i])
		_ = isNumber(number)
		i += 1
	}
	return Scan(input, i, n, arr, number, valName)
}
func appendIfExists(arr []string, number string, valName string) []string {
	if len(number) > 0 {
		arr = append(arr, number)
	}
	if len(valName) > 0 {
		arr = append(arr, valName)
	}
	return arr
}
func isReserved(input string) (string, bool) {

	for _, vi := range op {
		if string(input[0]) == string(vi) {
			return string(vi), true
		}
	}
	return "", false
}

func startWith(input string, v []string) (string, bool) {
	for _, vi := range v {
		if strings.HasPrefix(input, vi) {
			return vi, true
		}
	}
	return "", false

}
func isNumber(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return i
}
func isAlpha(c rune) bool {
	return (rune('a') <= c && c <= rune('z')) ||
		(rune('A') <= c && c <= rune('Z') || c == rune('_'))
}
func isAlphaNum(c rune) bool {
	res := isAlpha(c) || (rune('0') <= c && c <= rune('9'))
	return res
}
func getEscapeChar(c string) rune {
	switch c {
	case "a":
		return '\a'
	case "b":
		return '\b'
	case "t":
		return '\t'
	case "n":
		return '\n'
	case "v":
		return '\v'
	case "f":
		return '\f'
	case "r":
		return '\r'
	case "e":
		return 27
	case "0":
		return 0
	default:
		return rune(c[0])
	}

}
