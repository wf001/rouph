package lib

import (
	"fmt"
)

func LibPrint() {
	fmt.Println("")
}

func RouphLibPuts(name string, length int) {
	fmt.Printf("  push offset %s\n", name)
	fmt.Printf("  mov rsi, [rsp]\n")
	fmt.Printf("  mov rdx, %d\n", length)
	fmt.Printf("  mov rax, 1\n")
	fmt.Printf("  mov rdi, 1\n")
	fmt.Printf("  syscall\n")
	fmt.Printf("  pop rax\n")

	fmt.Printf("  push 0xa\n")
	fmt.Printf("  mov rsi, rsp\n")
	fmt.Printf("  mov rax, 1\n")
	fmt.Printf("  syscall\n")
	fmt.Printf("  pop rax\n")
	fmt.Printf("  ret\n")
}
