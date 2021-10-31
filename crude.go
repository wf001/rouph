package main

func main() {
	// tokenize
	head := TokenizeHandler()
	printToken(head)
	// parse
	prg := Program(head.Next)

	offset := 0
	for v := prg.Locals; v != nil; v = v.Next {
		offset += 8
		v.Offset = offset
	}
	prg.StackSize = offset

	Info("%s\n", "======== print node =========")
	n := prg.N
	for ; n != nil; n = n.Next {
		printNode(n)
		Info("%s\n", "## [ref next]")
	}
	Info("%s\n", "=================")
	// generate
	codegen(prg)

}
