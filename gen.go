package main

import (
	"fmt"
)

func genAddr(node *Node) {
	if node.Kind == ND_KIND_LVAR {
		offset := (rune(node.Name[0]) - rune('a') + 1) * 8
		fmt.Printf("  lea rax, [rbp-%d]\n", offset)
		fmt.Println("  push rax")
        return
	}
    panic("not an local value.")
}
func load(){
    fmt.Println("  pop rax")
    fmt.Println("  mov rax, [rax]")
    fmt.Println("  push rax")
}
func store(){
    fmt.Println("  pop rdi")
    fmt.Println("  pop rax")
    fmt.Println("  mov [rax], rdi")
    fmt.Println("  push rdi")
}

func gen(node *Node) {
	switch node.Kind {
	case ND_KIND_NUM:
		fmt.Printf("  push %d\n", node.Val)
		return
	case ND_KIND_EXPR_STMT:
		gen(node.Lhs)
		fmt.Println("  add rsp, 8")
		return
    case ND_KIND_LVAR:
        genAddr(node)
        load()
        return
    case ND_KIND_ASSIGN:
        genAddr(node.Lhs)
        gen(node.Rhs)
        store()
        return
	case ND_KIND_RETURN:
		gen(node.Lhs)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  jmp .Lreturn\n")
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
	Info("%s\n", "---------------------- instruction ---------------")
	Info("%s\n", "")
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")
    fmt.Println("  push rbp")
    fmt.Println("  mov rbp, rsp")
    fmt.Println("  sub rsp, 208")
	Info("%+v\n", node.Next)

	for ; node != nil; node = node.Next {
		gen(node)
	}
    fmt.Println(".Lreturn:")
    fmt.Println("  mov rsp, rbp")
    fmt.Println("  pop rbp")
	fmt.Println("  ret")
}
