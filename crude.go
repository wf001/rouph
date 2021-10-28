package main

import (
	"./crudego"
	"fmt"
)

func main() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println("main:")
	head := crudego.TokenizeHandler()
	crudego.Info("main head: %p\n", head)
	fmt.Println("  ret")
}
