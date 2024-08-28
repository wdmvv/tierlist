package main

import (
	"os"
)

func main() {
	ArgsParse()
	t := NewTierlist(os.Stdin, os.Stdout)
	if *Settings.Advanced {
		t.REPLAdvanced()
	} else {
		t.REPLBasic()
	}

}
