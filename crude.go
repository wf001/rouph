package main

func main() {
	head := TokenizeHandler()
	printToken(head)
	_, node := Program(head.Next)
    n := node
	Info("%s\n", "=================")
	for ; n != nil; n = n.Next {
        printNode(n)
        Info("%s\n", "## [rel next]")
    }
	Info("%s\n", "=================")
	codegen(node)

}
