package main

type TypeKind int

const (
	TY_CHAR TypeKind = iota + 1 // +
	TY_INT
	TY_PTR
	TY_ARRAY
)

type Type struct {
	Kind      TypeKind
	Base      *Type
	ArraySize int
}

func newType(kind TypeKind) *Type {
	ty := new(Type)
	ty.Kind = kind
	return ty
}

func intType() *Type {
	return newType(TY_INT)
}
func charType() *Type {
	return newType(TY_CHAR)
}
func pointerTo(base *Type) *Type {
	t := newType(TY_PTR)
	t.Base = base
	return t
}
func arrayOf(base *Type, size int) *Type {
	ty := newType(TY_ARRAY)
	ty.Base = base
	ty.ArraySize = size
	return ty
}
func sizeOf(ty *Type) int {
	switch ty.Kind {
	case TY_INT, TY_PTR:
		return 8
	case TY_CHAR:
		return 1
	default:
		if ty.Kind != TY_ARRAY {
			panic("type is allowed only array.")
		}
		return sizeOf(ty.Base) * ty.ArraySize
	}
}

func visit(node *Node) {
	if node == nil {
		return
	}
	visit(node.Lhs)
	visit(node.Rhs)
	visit(node.Cond)
	visit(node.Then)
	visit(node.Else)
	visit(node.Init)
	visit(node.Inc)

	for b := node.Body; b != nil; b = b.Next {
		visit(b)
	}
	for a := node.Args; a != nil; a = a.Next {
		visit(a)
	}

	switch node.Kind {
	case ND_KIND_MUL,
		ND_KIND_DIV,
		ND_KIND_EQ,
		ND_KIND_NE,
		ND_KIND_LE,
		ND_KIND_LT,
		ND_KIND_GE,
		ND_KIND_GT,
		ND_KIND_FUNCALL,
		ND_KIND_NUM:
		node.Ty = intType()
		return
	case ND_KIND_VAR:
		node.Ty = node.Var.Ty
		return
	case ND_KIND_ADD:
		if node.Rhs.Ty.Base != nil {
			tmp := node.Lhs
			node.Lhs = node.Rhs
			node.Rhs = tmp
		}
		if node.Rhs.Ty.Base != nil {
			panic("invalid token")
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_SUB:
		if node.Rhs.Ty.Base != nil {
			panic("invalid token")
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_ASSIGN:
		node.Ty = node.Lhs.Ty
		return
	case ND_KIND_ADDR:
		node.Ty = pointerTo(node.Lhs.Ty)
		if node.Lhs.Ty.Kind == TY_ARRAY {
			node.Ty = pointerTo(node.Lhs.Ty.Base)
		} else {
			node.Ty = pointerTo(node.Lhs.Ty)
		}
		return
	case ND_KIND_DEREF:
		if node.Lhs.Ty.Base == nil {
			panic("Invalid pointer dereference")
		}
		node.Ty = node.Lhs.Ty.Base
		return
	case ND_KIND_SIZEOF:
		node.Kind = ND_KIND_NUM
		node.Ty = intType()
		node.Val = sizeOf(node.Lhs.Ty)
		node.Lhs = nil
	}
}
func addType(prg *Prog) {
	for fn := prg.Fns; fn != nil; fn = fn.Next {
		for n := fn.N; n != nil; n = n.Next {
			visit(n)
		}
	}
}
