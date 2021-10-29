package main

import (
	"fmt"
)
func main() {
	head := TokenizeHandler()
	printToken(head)
    _, node := Expr(head.Next)

	Info("%s\n", "=================")
	printNode(node)
	Info("%s\n", "=================")

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")
    gen(node)
	fmt.Println("  pop rax")
	fmt.Println("  ret")
}
