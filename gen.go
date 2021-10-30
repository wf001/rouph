package main

import (
	"fmt"
)

func gen(node *Node) {
	if node.Kind == ND_KIND_NUM {
		fmt.Printf("  push %d\n", node.Val)
		return
	}
	gen(node.Lhs)
	gen(node.Rhs)
	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	switch node.Kind {
	case ND_KIND_ADD:
		fmt.Printf("  add rax, rdi\n")
		break
	case ND_KIND_SUB:
		fmt.Printf("  sub rax, rdi\n")
		break
	case ND_KIND_MUL:
		fmt.Printf("  imul rax, rdi\n")
		break
	case ND_KIND_DIV:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
		break
	case ND_KIND_EQ:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  sete al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	case ND_KIND_NE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setne al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	case ND_KIND_LT:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setl al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	case ND_KIND_LE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setle al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	case ND_KIND_GT:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setg al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	case ND_KIND_GE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  setge al\n")
		fmt.Printf("  movzb rax, al\n")
		break
	}
	fmt.Println("  push rax")
}
func codegen(node *Node) {
	Info("%s\n","---------------------- instruction ---------------")
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")

	for ; node != nil; node = node.Next {
		gen(node)
		fmt.Println("  pop rax")
	}
	fmt.Println("  ret")
}
