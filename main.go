package main

import (
	"os"
)

func main() {
	ArgsParse()
	t := NewTierlist(os.Stdin, os.Stdout)
	t.REPL()
}
