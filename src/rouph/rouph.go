package main

import (
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	flag.Parse()
	file := flag.Arg(0)
	out, err := exec.Command("sh", "-c", "./bin/rouphc "+file+" > test.s").CombinedOutput()
	if err != nil {
		fmt.Printf("CombineOutput: %s, Error: %v\n", out, err)
	}

	out, err = exec.Command("gcc", "-static", "-o", "test", "test.s").CombinedOutput()
	if err != nil {
		fmt.Printf("CombineOutput: %s, Error: %v\n", out, err)
	}

	_, err = exec.Command("rm", "-f", "test.s").CombinedOutput()
	if err != nil {
		fmt.Printf("CombineOutput: %s, Error: %v\n", out, err)
	}
}
