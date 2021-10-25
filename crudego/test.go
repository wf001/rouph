package crudego

import (
	"flag"
	"fmt"
	"regexp"
)

type TokKind int

const (
	TK_KIND_RESERVED TokKind = iota + 1
	TK_KIND_NUM
	TK_KIND_EOF
)

type Token struct {
	kind TokKind
	next *Token
	val  string
}

func newToken(kind TokKind, cur *Token, val string) *Token {
	fmt.Printf("kind: %d, cur: %+v, val: %s\n", kind, cur, val)
	tok := new(Token)
	tok.kind = kind
	tok.val = val
	cur.next = tok
	return tok
}
func printToken(tok *Token) {
	fmt.Println("==========================")
	for ; tok != nil; tok = tok.next {
        fmt.Printf("tok:%p, %+v\n", tok, tok)
	}
	fmt.Println("==========================")
}
func Hoge() {
	flag.Parse()
	arg := flag.Arg(0)
	reg := "[+-]"
	arg_arr := regexp.MustCompile(reg).Split(arg, -1)
	cur_len := len(arg_arr[0])

	head := new(Token)
	head.next = nil
	head.kind = -1
	cur := head

	cur = newToken(TK_KIND_NUM, cur, arg_arr[0])

	for _, s := range arg_arr[1:] {
		if string(arg[cur_len]) == "+" || string(arg[cur_len]) == "-" {
			cur = newToken(TK_KIND_RESERVED, cur, string(arg[cur_len]))
			cur = newToken(TK_KIND_NUM, cur, s)
			cur_len += len(s) + 1
		}
	}
	cur = newToken(TK_KIND_EOF, cur, "")
    printToken(head)
}
