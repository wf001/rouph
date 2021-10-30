package main

func main() {
	head := TokenizeHandler()
	printToken(head)
	_, node := Program(head.Next)
    n := node
	Info("%s\n", "======== print node =========")
	for ; n != nil; n = n.Next {
        printNode(n)
        Info("%s\n", "## [ref next]")
    }
	Info("%s\n", "=================")
	codegen(node)

}
