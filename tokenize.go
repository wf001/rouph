package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var reserve_sig = []string{"return", "if", "else", "for", "int", "sizeof", "char"}
var eq_rel_op = []string{"==", "!=", ">=", "=<"}
var op = "+-*/()><;={},&[]"

type TokKind int

const (
	TK_KIND_RESERVED TokKind = iota + 1
	TK_IDENT
	TK_STR
	TK_KIND_NUM
	TK_KIND_EOF
)

type Token struct {
	Kind TokKind
	Next *Token
	Val  string
	Str  string
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
func TokenizeHandler() *Token {
	flag.Parse()
	arg := flag.Arg(0)
	// space trim
	Info("arg:%s\n", arg)
	// gen num arr
	var arg_arr []string
	arg_arr = Scan(arg, 0, len(arg), arg_arr, "", "")

	head := new(Token)
	head.Next = nil
	head.Kind = -1
	cur := head

	//tokenize
	Info("arg_arr:%#v\n", arg_arr)
	for _, s := range arg_arr {
		if _, res := isReserved(s); res {
			cur = newToken(TK_KIND_RESERVED, cur, s)
		} else if s == "return" {
			cur = newToken(TK_KIND_RESERVED, cur, s)
		} else if isAlpha(rune(s[0])) {
			cur = newToken(TK_IDENT, cur, s)
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
	// if input[i] is space, it skipped
	if string(input[i]) == " " {
		arr = appendIfExists(arr, number, valName)
		number, valName = "", ""
		i += 1
		// if input[i] is reserved signature
	} else if sig, res := startWith(string(input[i:]), reserve_sig); res &&
		len(valName) == 0 &&
		!isAlphaNum(rune(input[i+len(sig)])) {
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
