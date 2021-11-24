package main

import (
	"bufio"
	"flag"
	"os"
)

func alignTo(n int, align int) int {
	return (n + align - 1) & (^(align - 1))
}

func main() {
	s := flag.String("i", "", "")
	flag.Parse()
	var arg string
    Info("%s\n", *s)

	if *s != "" {
		arg = *s
	} else {
        file := flag.Arg(0)
		data, _ := os.Open(file)
		defer data.Close()
		scanner := bufio.NewScanner(data)
		var input_arr = ""
		for scanner.Scan() {
			input_arr += scanner.Text()
		}
        arg = input_arr
	}
	Info("arg:%s\n", arg)
	// tokenize
	head := TokenizeHandler(arg)
	printToken(head)
	// parse
	prg := Program(head.Next)
	Info("%+v\n", prg)
	addType(prg)
	Info("%+v\n", prg)

	for fn := prg.Fns; fn != nil; fn = fn.Next {
		offset := 0

		for vl := fn.Locals; vl != nil; vl = vl.Next {
			v := vl.V
			offset += sizeOf(v.Ty)
			v.Offset = offset
		}
		fn.StackSize = alignTo(offset, 8)
	}

	Info("%s\n", "======== print node =========")
	for fn := prg.Fns; fn != nil; fn = fn.Next {
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
