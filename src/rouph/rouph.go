package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)
	file := flag.Arg(1)
	fs := strings.Split(file, "/")
	exec_file := strings.Split(fs[len(fs)-1], ".")[0]
	if cmd == "build" {
		rouphBuild(file, exec_file)
	} else if cmd == "run" {
		rouphBuild(file, exec_file)
		rouphRun(exec_file)
	} else {
		rouphHelp()
	}

}

func rouphBuild(file string, exec_file string) {
	out, err := exec.Command("sh", "-c", "./bin/rouphc "+file+" > test.s").CombinedOutput()
	if err != nil {
		fmt.Printf("1rouph: %s, Error: %v\n", out, err)
	}

	out, err = exec.Command("gcc", "-static", "-o", exec_file, "test.s").CombinedOutput()
	if err != nil {
		fmt.Printf("rouph: Error: %v", err)
	}
	_, err = exec.Command("rm", "-f", "test.s").CombinedOutput()
	if err != nil {
		fmt.Printf("rouph: %s, Error: %v\n", out, err)
	}

}
func rouphRun(exec_file string) {
	out, err := exec.Command("./" + exec_file).CombinedOutput()
	fmt.Printf("%s\n", out)
	if err != nil {
		fmt.Printf("rouph: Error: %v\n", err)
	}
	_, err = exec.Command("rm", "-f", exec_file).CombinedOutput()
	if err != nil {
		fmt.Printf("rouph: %s, Error: %v\n", out, err)
	}

}

func rouphHelp() {
	fmt.Printf("this is help.\n")
}
