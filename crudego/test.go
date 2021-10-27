package crudego

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
)

var DEBUG bool = false

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

func printToken(tok *Token) {
	if DEBUG {

		fmt.Println("==========================")
		for ; tok != nil; tok = tok.next {
			info("tok %p, ", tok)
			info("%+v\n", tok)
		}
		fmt.Println("==========================")
	}
}
func info(s string, v interface{}) {
	if DEBUG {
		fmt.Printf(s, v)
	}
}

func newToken(kind TokKind, cur *Token, val string) *Token {
	info("kind: %d, ", kind)
	info("cur: %+v, ", cur)
	info("val: %s\n", val)
	tok := new(Token)
	tok.kind = kind
	tok.val = val
	cur.next = tok
	return tok
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

	fmt.Printf("  mov rax, %d\n", 0)
	for e := head.next; e.next != nil; e = e.next {
		info("e.val: %s\n", e.val)
		if e.val == "+" {
			e = e.next
			fmt.Printf("  add rax, %s\n", e.val)
			continue
		} else if e.val == "-" {
			e = e.next
			fmt.Printf("  sub rax, %s\n", e.val)
			continue
		}
		i, err := strconv.Atoi(e.val)
		if err != nil {
			panic("invalid")
		}
		fmt.Printf("  add rax, %d\n", i)
	}
}
