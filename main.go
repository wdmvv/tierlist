package main

import (
	"os"
)

func main() {
	ArgsParse()
	t := NewTierlist(os.Stdin, os.Stdout)
	// if *Settings.Preset != 0{
	// 	will add in the future 'updates'
	// }
	if *Settings.Advanced {
		t.REPLAdvanced()
	} else {
		t.REPLBasic()
	}

}
