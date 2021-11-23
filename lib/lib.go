package lib

import (
	"fmt"
)

func LibPrint() {
	fmt.Println("")
}
func StdlibHandler() {
	libPuts()
}

func libPuts() {
	fmt.Printf(".puts:\n")
	fmt.Printf("  mov rdx, 1\n")
	fmt.Printf("  mov rax, 1\n")
	fmt.Printf("  mov rdi, 1\n")
	fmt.Printf("  syscall\n")
	fmt.Printf("  sub r8, 1\n")
	fmt.Printf("  add rsi, 1\n")
	fmt.Printf("  cmp r8, 0\n")
	fmt.Printf("  jnz .puts\n")

	fmt.Printf("  push 0xa\n")
	fmt.Printf("  mov rsi, rsp\n")
	fmt.Printf("  mov rax, 1\n")
	fmt.Printf("  syscall\n")
	fmt.Printf("  pop rax\n")
    fmt.Printf("  mov rax, 1\n")
	fmt.Printf("  ret\n")
}
