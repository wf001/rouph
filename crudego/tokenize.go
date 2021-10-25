package crudego

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)


func TokenizeHandler() {
	flag.Parse()
	if len(flag.Args()) > 1 {
		fmt.Fprintln(os.Stderr, "Invalid args")
		os.Exit(-1)
	}

	arg := flag.Arg(0)
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	reg := "[+-]"
	arg_arr := regexp.MustCompile(reg).Split(arg, -1)
	cur_len := len(arg_arr[0])

	val, _ := strconv.Atoi(arg_arr[0])
	fmt.Printf("  mov rax, %d\n", val)

	//Parse
	for _, s := range arg_arr[1:] {
		val, _ = strconv.Atoi(s)
		if string(arg[cur_len]) == "+" {
			fmt.Printf("  add rax, %d\n", val)
		} else if string(arg[cur_len]) == "-" {
			fmt.Printf("  sub rax, %d\n", val)
		} else {
			return
		}
		cur_len += len(s) + 1
	}

	fmt.Println("  ret")
}
