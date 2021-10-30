package main

func main() {
	head := TokenizeHandler()
	printToken(head)
	_, node := Expr(head.Next)

	Info("%s\n", "=================")
	printNode(node)
	Info("%s\n", "=================")
	codegen(node)

}
