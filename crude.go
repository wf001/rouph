package main // mainパッケージであることを宣言

import (
	"flag"
	"fmt"
    "strconv"
)

func main() { // 最初に実行されるmain()関数を定義
    flag.Parse()
    var arg int
    arg,_ = strconv.Atoi(flag.Arg(0))
    fmt.Println(".intel_syntax noprefix")
    fmt.Println(".global main")
    fmt.Println("main:")
    fmt.Printf("  mov rax, %d\n", arg)
    fmt.Println("  ret")
}
