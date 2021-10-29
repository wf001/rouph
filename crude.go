package main

import (
	"fmt"
	"crude-lang-go/lib"
)

func main() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")
	head := TokenizeHandler()
    _ = Expr(head.Next)
	Info("main head: %p\n", head)
	fmt.Println("  ret")
    lib.LibPrint()
}
