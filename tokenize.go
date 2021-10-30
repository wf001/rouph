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

type TokKind int

const (
	TK_KIND_RESERVED TokKind = iota + 1
	TK_KIND_NUM
	TK_KIND_EOF
)

type Token struct {
	Kind TokKind
	Next *Token
	Val  string
	str  string
}

/*
Token Func
*/
func printToken(tok *Token) {
	if DEBUG {
		fmt.Println("==========================")
		for ; tok != nil; tok = tok.Next {
			Info("tok %p\n", tok)
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
	tok.Val = val
	cur.Next = tok
	return tok
}
func TokenizeHandler() *Token {
	flag.Parse()
	arg := flag.Arg(0)
	// space trim
	arg = strings.Replace(arg, " ", "", -1)
	Info("arg:%s\n", arg)
	// gen num arr
	var arg_arr []string
	arg_arr = Scan(arg, 0, len(arg), arg_arr, "")

	head := new(Token)
	head.Next = nil
	head.Kind = -1
	cur := head

	//tokenize
	for _, s := range arg_arr {
		if _, res := strChr(s); res {
			cur = newToken(TK_KIND_RESERVED, cur, s)
		} else {
			cur = newToken(TK_KIND_NUM, cur, s)
		}
	}
	cur = newToken(TK_KIND_EOF, cur, "")

	return head
}

func Scan(input string, i int, n int, arr []string, number string) []string {
	if i == n {
		arr = appendIfExists(arr, number)
		return arr
	}
	// if input[i] equ or rel
	if sig, res := isEquOrRel(string(input[i:])); res {
		arr = appendIfExists(arr, number)
		arr, number = append(arr, sig), ""
		i += 2
		// if input[i] is other reserved
	} else if sig, res := strChr(string(input[i])); res {
		arr = appendIfExists(arr, number)
		arr, number = append(arr, sig), ""
		i += 1
		// if input[i] is number
	} else {
		number += string(input[i])
		_ = isNumber(number)
		i += 1
	}
	return Scan(input, i, n, arr, number)
}
func appendIfExists(arr []string, number string) []string {
	if len(number) > 0 {
		arr = append(arr, number)
	}
	return arr
}
func strChr(input string) (string, bool) {
	v := "+-*/()><"

	for _, vi := range v {
		if string(input[0]) == string(vi) {
			return string(vi), true
		}
	}
	return "", false
}

func isEquOrRel(input string) (string, bool) {
	v := []string{"==", "!=", ">=", "=<"}

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
