package main

import (
	"fmt"
	"rouphc/lib"
)

var argReg1 = []string{"dil", "sil", "dl", "cl", "r8b", "r9b"}
var argReg8 = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var labelSeq = 0
var funcName string

func genAddr(node *Node) {
	switch node.Kind {
	case ND_KIND_VAR:
		v := node.Var
		if v.isLocal {
			fmt.Printf("  lea rax, [rbp-%d]\n", v.Offset)
			fmt.Println("  push rax")
		} else {
			fmt.Printf("  push offset %s\n", v.Name)
		}
		return
	case ND_KIND_DEREF:
		gen(node.Lhs)
		return
	}
	panic("not an local value.")
}
func genLval(node *Node) {
	if node.Ty.Kind == TY_ARRAY {
		panic("not a local value.")
	}
	genAddr(node)
}
func load(ty *Type) {
	fmt.Println("  pop rax")
	if sizeOf(ty) == 1 {
		fmt.Println("  movsx rax, byte ptr [rax]")
	} else {
		fmt.Println("  mov rax, [rax]")
	}
	fmt.Println("  push rax")
}
func store(ty *Type) {
	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")
	if sizeOf(ty) == 1 {
		fmt.Println("  mov [rax], dil")
	} else {
		fmt.Println("  mov [rax], rdi")
	}
	fmt.Println("  push rdi")
}

func gen(node *Node) {
	switch node.Kind {
	case ND_KIND_NULL:
		return
	case ND_KIND_NUM:
		fmt.Printf("  push %d\n", node.Val)
		return
	case ND_KIND_EXPR_STMT:
		gen(node.Lhs)
		fmt.Println("  add rsp, 8")
		return
	case ND_KIND_VAR:
		genAddr(node)
		if node.Ty.Kind != TY_ARRAY {
			load(node.Ty)
		}
		return
	case ND_KIND_ASSIGN:
		// push local val address
		genLval(node.Lhs)
		// push right side val
		gen(node.Rhs)
		store(node.Ty)
		return
	case ND_KIND_ADDR:
		genAddr(node.Lhs)
		return
	case ND_KIND_DEREF:
		gen(node.Lhs)
		if node.Ty.Kind != TY_ARRAY {
			load(node.Ty)
		}
		return
	case ND_KIND_IF:
		var seq = labelSeq
		labelSeq++
		if node.Else != nil {
			gen(node.Cond)
			fmt.Println("  pop rax")
			fmt.Println("  cmp rax, 0")
			fmt.Printf("  je .Lelse%d\n", seq)
			gen(node.Then)
			fmt.Printf("  jmp .Lend%d\n", seq)
			fmt.Printf(".Lelse%d:\n", seq)
			gen(node.Else)
			fmt.Printf(".Lend%d:\n", seq)
		} else {
			gen(node.Cond)
			fmt.Println("  pop rax")
			fmt.Println("  cmp rax,0")
			fmt.Printf("  je .Lend%d\n", seq)
			gen(node.Then)
			fmt.Printf(".Lend%d:\n", seq)
		}
		return
	case ND_KIND_FOR:
		seq := labelSeq
		labelSeq++
		if node.Init != nil {
			gen(node.Init)
		}
		fmt.Printf(".Lbegin%d:\n", seq)
		if node.Cond != nil {
			gen(node.Cond)
			fmt.Println("  pop rax")
			fmt.Println("  cmp rax, 0")
			fmt.Printf("  je .Lend%d\n", seq)
		}
		gen(node.Then)
		if node.Inc != nil {
			gen(node.Inc)
		}
		fmt.Printf("  jmp .Lbegin%d\n", seq)
		fmt.Printf(".Lend%d:\n", seq)
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
			fmt.Printf("  pop %s\n", argReg8[i])
		}

		seq := labelSeq
		labelSeq += 1
		fmt.Println("  mov rax, rsp")
		fmt.Println("  and rax, 15")
		fmt.Printf("  jnz .Lcall%d\n", seq)
		fmt.Println("  mov rax, 0")
		fmt.Printf("  call %s\n", node.Func)
		fmt.Printf("  jmp .Lend%d\n", seq)

		fmt.Printf(".Lcall%d:\n", seq)
		fmt.Println("  sub rsp, 8")
		fmt.Println("  mov rax, 0")
		fmt.Printf("  call %s\n", node.Func)
		fmt.Println("  add rsp, 8")

		fmt.Printf(".Lend%d:\n", seq)
		fmt.Printf("  push rax\n")
		return
	case ND_KIND_RETURN:
		gen(node.Lhs)
		fmt.Printf("  pop rax\n")
		fmt.Printf("  jmp .Lreturn.%s\n", funcName)
		return
	case ND_KIND_STDLIB:
		Info("stdlib %+v\n", node)
		Info("stdlib %+v\n", node.Lhs.Var)
		Info("stdlib %+v\n", node.Lhs.Var.Ty)
		Info("stdlib %+v\n", node.Lhs.Var.Ty.Base)
		if node.Lhs.Var != nil &&
			node.Lhs.Var.isLocal &&
			(node.Lhs.Var.Ty.Kind != TY_PTR || node.Lhs.Var.Ty.Base.Kind != TY_CHAR) {

			fmt.Printf("  lea r10, [rsp-8]\n")
			fmt.Printf("  mov rax, [rbp-%d]\n", node.Lhs.Var.Offset)
			fmt.Printf("  mov rbx, 10\n")
			fmt.Printf("  mov r8, %d\n", 0)
			//Type identifier 0 = int
			fmt.Printf("  mov r9, %d\n", 0)
			fmt.Printf("  mov r12, rsp\n")
			fmt.Printf("  call .%s\n", node.Lhs.Func)
			fmt.Printf("  mov rsp, r12\n")
			fmt.Printf("  sub rsp, 8\n")
		} else if node.Lhs.Var != nil &&
			node.Lhs.Var.isLocal &&
			node.Lhs.Var.Ty.Kind == TY_PTR && node.Lhs.Var.Ty.Base.Kind == TY_CHAR {

			fmt.Printf("  lea r10, [rsp-8]\n")
			fmt.Printf("  mov rax, [rbp-%d]\n", node.Lhs.Var.Offset)
			fmt.Printf("  mov rbx, 10\n")
			fmt.Printf("  mov r8, %d\n", 0)
			//Type identifier 2 = local string
			fmt.Printf("  mov r9, %d\n", 2)
			fmt.Printf("  mov r12, rsp\n")
			fmt.Printf("  mov rsi, [rbp-%d]\n", node.Lhs.Var.Offset)
			fmt.Printf("  call .%s\n", node.Lhs.Func)
			fmt.Printf("  mov rsp, r12\n")
			fmt.Printf("  sub rsp, 8\n")
		} else {
			fmt.Printf("  push offset %s\n", node.Lhs.Var.Name)
			fmt.Printf("  mov r8, %d\n", node.Lhs.Var.ContLen-1)
			fmt.Printf("  mov rsi, [rsp]\n")
			//Type identifier 1 = string
			fmt.Printf("  mov r9, %d\n", 1)
			fmt.Printf("  mov r12, rsp\n")
			fmt.Printf("  call .%s\n", node.Lhs.Func)
		}
		return
	}
	gen(node.Lhs)
	gen(node.Rhs)
	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	switch node.Kind {
	case ND_KIND_ADD:
		if node.Ty.Base != nil {
			fmt.Printf("  imul rdi, %d\n", sizeOf(node.Ty.Base))
		}
		fmt.Printf("  add rax, rdi\n")
		break
	case ND_KIND_SUB:
		if node.Ty.Kind == TY_PTR {
			fmt.Printf("  imul rdi, %d\n", sizeOf(node.Ty.Base))
		}
		fmt.Printf("  sub rax, rdi\n")
		break
	case ND_KIND_MUL:
		fmt.Printf("  imul rax, rdi\n")
		break
	case ND_KIND_DIV:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
		break
	case ND_KIND_REM:
		fmt.Printf("  mov rdx, 0\n")
		fmt.Printf("  mov rbx, rdi\n")
		fmt.Printf("  div rbx\n")
		fmt.Printf("  mov rax, rdx\n")
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
func emitData(prg *Prog) {
	fmt.Printf(".data\n")
	for vl := prg.Globals; vl != nil; vl = vl.Next {
		v := vl.V
		fmt.Printf("%s:\n", v.Name)
		if v.ContLen == 0 {
			fmt.Printf("  .zero %d\n", sizeOf(v.Ty))
			continue
		}
		for i := 0; i < v.ContLen; i++ {
			fmt.Printf("  .byte %d\n", v.Contents[i])
		}
	}
}

func loadArg(v *Var, idx int) {
	sz := sizeOf(v.Ty)
	if sz == 1 {
		fmt.Printf("  mov [rbp-%d], %s\n", v.Offset, argReg1[idx])
	} else if sz == 8 {
		fmt.Printf("  mov [rbp-%d], %s\n", v.Offset, argReg8[idx])
	} else {
		panic("type size must be 1 or 8.")
	}
}

func emitText(prg *Prog) {
	fmt.Printf(".text\n")

	for fn := prg.Fns; fn != nil; fn = fn.Next {
		fmt.Printf(".global %s\n", fn.Name)
		fmt.Printf("%s:\n", fn.Name)
		funcName = fn.Name

		fmt.Printf("  push rbp\n")
		fmt.Printf("  mov rbp, rsp\n")
		fmt.Printf("  sub rsp, %d\n", fn.StackSize)

		i := 0
		for vl := fn.Params; vl != nil; vl = vl.Next {
			loadArg(vl.V, i)
			i += 1
		}

		for node := fn.N; node != nil; node = node.Next {
			gen(node)
		}
		fmt.Printf(".Lreturn.%s:\n", funcName)
		fmt.Println("  mov rsp, rbp")
		fmt.Println("  pop rbp")
		fmt.Println("  ret")

	}

}
func codegen(prg *Prog) {
	Info("%s\n", "---------------------- instruction ---------------")
	Info("%s\n", "")
	fmt.Println(".intel_syntax noprefix")
	emitData(prg)
	emitText(prg)
	lib.StdlibHandler()
}
