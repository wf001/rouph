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
	fmt.Printf(".put:\n")
	fmt.Printf("  cmp r9, 1\n")
	fmt.Printf("  je .puts\n")
	fmt.Printf("  jne .store\n")
	fmt.Printf("  ret\n")

	fmt.Printf(".store:\n")
	fmt.Printf("  mov rdx, 0\n")
	fmt.Printf("  div rbx\n")
	fmt.Printf("  add rdx, 0x30\n")
	fmt.Printf("  push rdx\n")
	fmt.Printf("  inc r8\n")
	fmt.Printf("  cmp rax, 0\n")
	fmt.Printf("  jnz .store\n")
	fmt.Printf("  lea rsi, [rsp]\n")
	fmt.Printf("  call .puti\n")
	fmt.Printf("  lea rsp, [rbp]\n")
	fmt.Printf("  ret\n")

	fmt.Printf(".puti:\n")
	fmt.Printf("  mov rdx, 1\n")
	fmt.Printf("  mov rax, 1\n")
	fmt.Printf("  mov rdi, 1\n")
	fmt.Printf("  syscall\n")
	fmt.Printf("  sub r8, 1\n")
	fmt.Printf("  add rsi, 8\n")
	fmt.Printf("  cmp r8, 0\n")
	fmt.Printf("  jnz .puti\n")
	fmt.Printf("  push 0xa\n")
	fmt.Printf("  mov rsi, rsp\n")
	fmt.Printf("  mov rax, 1\n")
	fmt.Printf("  syscall\n")
	fmt.Printf("  pop rax\n")
	fmt.Printf("  mov rax, 1\n")
	fmt.Printf("  ret\n")

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
