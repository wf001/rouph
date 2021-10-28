package crudego

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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
}

func printToken(tok *Token) {
	if DEBUG {

		fmt.Println("==========================")
		for ; tok != nil; tok = tok.Next {
			Info("tok %p, ", tok)
			Info("%+v\n", tok)
		}
		fmt.Println("==========================")
	}
}
func Info(s string, v interface{}) {
	if DEBUG {
		fmt.Printf(s, v)
	}
}

func newToken(kind TokKind, cur *Token, val string) *Token {
	Info("kind: %d, ", kind)
	Info("cur: %+v, ", cur)
	Info("val: %s\n", val)
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
	// gen num arr
	reg := "[+-]"
	arg_arr := regexp.MustCompile(reg).Split(arg, -1)
	cur_len := len(arg_arr[0])

	head := new(Token)
	head.Next = nil
	head.Kind = -1
	cur := head

	cur = newToken(TK_KIND_NUM, cur, arg_arr[0])

	//tokenize
	for _, s := range arg_arr[1:] {
		if string(arg[cur_len]) == "+" || string(arg[cur_len]) == "-" {
			cur = newToken(TK_KIND_RESERVED, cur, string(arg[cur_len]))
			cur = newToken(TK_KIND_NUM, cur, s)
			cur_len += len(s) + 1
		}
	}
	cur = newToken(TK_KIND_EOF, cur, "")
	fmt.Printf("  mov rax, %d\n", 0)
	//parsing
	for e := head.Next; e.Next != nil; e = e.Next {
		Info("e.val: '%s'\n", e.Val)
		var op string
		if e.Val == "+" {
			e = e.Next
			op = "  add rax "
		} else if e.Val == "-" {
			op = "  sub rax "
			e = e.Next
		} else {
			op = "  add rax "
		}
		i := isNumber(e.Val)
		fmt.Printf("  %s, %d\n", op, i)
	}

	return head
}

func isNumber(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return i
}
