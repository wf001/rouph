package main

import (
	"fmt"
)

var labelSeq = 0
var argReg = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func genAddr(node *Node) {
	if node.Kind == ND_KIND_VAR {
		fmt.Printf("  lea rax, [rbp-%d]\n", node.Var.Offset)
		fmt.Println("  push rax")
		return
	}
	panic("not an local value.")
}
func load() {
	fmt.Println("  pop rax")
	fmt.Println("  mov rax, [rax]")
	fmt.Println("  push rax")
}
func store() {
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
	case ND_KIND_VAR:
		genAddr(node)
		load()
		return
	case ND_KIND_ASSIGN:
		// push local val address
		genAddr(node.Lhs)
		// push right side val
		gen(node.Rhs)
		store()
		return
	case ND_KIND_IF:
		var seq = labelSeq
		labelSeq++
		if node.Else != nil {
			gen(node.Cond)
			fmt.Println("  pop rax")
			fmt.Println("  cmp rax, 0")
			fmt.Printf("  je Lelse%d\n", seq)
			gen(node.Then)
			fmt.Printf("  je Lend%d\n", seq)
			fmt.Printf("Lelse%d:\n", seq)
			gen(node.Else)
			fmt.Printf("Lend%d:\n", seq)
		} else {
			gen(node.Cond)
			fmt.Println("  pop rax")
			fmt.Println("  cmp rax,0")
			fmt.Printf("  je Lend%d\n", seq)
			gen(node.Then)
			fmt.Printf("Lend%d:\n", seq)
		}
		return
	case ND_KIND_FOR:
		seq := labelSeq
		labelSeq++
		if node.Init != nil {
			gen(node.Init)
		}
		fmt.Printf("Lbegin%d:\n", seq)
		if node.Cond != nil {
			gen(node.Cond)
			fmt.Println("  pop rax")
			fmt.Println("  cmp rax, 0")
			fmt.Printf("  je Lend%d\n", seq)
		}
		gen(node.Then)
		if node.Inc != nil {
			gen(node.Inc)
		}
		fmt.Printf("  jmp Lbegin%d\n", seq)
		fmt.Printf("Lend%d:\n", seq)
		return
	case ND_KIND_BLOCK:
		for n := node.Body; n != nil; n = n.Next {
			gen(n)
		}
		return
	case ND_KIND_FUNCALL:
		nargs := 0
		for arg := node.Args; arg != nil; arg = arg.Next {
			gen(arg)
			nargs += 1
		}
		for i := nargs - 1; i >= 0; i -= 1 {
			fmt.Printf("  pop %s\n", argReg[i])
		}

		seq := labelSeq
		labelSeq += 1
		fmt.Println("  mov rax, rsp")
		fmt.Println("  and rax, 15")
		fmt.Printf("  jnz Lcall%d\n", seq)
		fmt.Println("  mov rax, 0")
		fmt.Printf("  call %s\n", node.Func)
		fmt.Printf("  jmp Lend%d\n", seq)

		fmt.Printf("Lcall%d:\n", seq)
		fmt.Println("  sub rsp, 8")
		fmt.Println("  mov rax, 0")
		fmt.Printf("  call %s\n", node.Func)
		fmt.Println("  add rsp, 8")

		fmt.Printf(" Lend%d:\n", seq)
		fmt.Printf("  push rax\n")
		return
	case ND_KIND_RETURN:
		gen(node.Lhs)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  jmp Lreturn\n")
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
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  sete al\n")
		break
	case ND_KIND_NE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  setne al\n")
		break
	case ND_KIND_LT:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  setl al\n")
		break
	case ND_KIND_LE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  setle al\n")
		break
	case ND_KIND_GT:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  setg al\n")
		break
	case ND_KIND_GE:
		fmt.Printf("  cmp rax, rdi\n")
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  setge al\n")
		break
	}
	fmt.Println("  push rax")
}
func codegen(prg *Prg) {
	Info("%s\n", "---------------------- instruction ---------------")
	Info("%s\n", "")
	fmt.Println("global _start")
	fmt.Println("_start:")
	fmt.Println("  push rbp")
	fmt.Println("  mov rbp, rsp")
	fmt.Printf("  sub rsp, %d\n", prg.StackSize)

	for node := prg.N; node != nil; node = node.Next {
		gen(node)
	}
	fmt.Println("Lreturn:")
	fmt.Println("  mov rbx, rax")
	fmt.Println("  mov rax, 1")
	fmt.Println("  int 0x80")
	fmt.Println("  ret")
}
