package main

func main() {
	// tokenize
	head := TokenizeHandler()
	printToken(head)
	// parse
	prg := Program(head.Next)
	Info("%+v\n", prg)
	addType(prg)
	Info("%+v\n", prg)

	for fn := prg; fn != nil; fn = fn.Next {
		offset := 0

		for vl := fn.Locals; vl != nil; vl = vl.Next {
			offset += 8
			vl.V.Offset = offset
		}
		fn.StackSize = offset
	}

	Info("%s\n", "======== print node =========")
	for fn := prg; fn != nil; fn = fn.Next {
		n := fn.N
		for ; n != nil; n = n.Next {
			printNode(n)
			Info("%s\n", "## [ref next]")
		}
	}
	Info("%s\n", "=================")
	// generate
	codegen(prg)

}
