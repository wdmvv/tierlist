package main

import (
	"os"
)

func main() {

	t := NewTierlist(os.Stdin, os.Stdout)
	t.REPL()
}
