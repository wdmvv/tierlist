package main

import (
	"os"
	"tierlist/cmd"
	"tierlist/pkg/tlist"
)

func main() {

	cmd.ArgsParse()
	t := tlist.NewTierlist(os.Stdin, os.Stdout)
	// if *Settings.Preset != 0{
	// 	will add in the future 'updates'
	// }
	if *cmd.Settings.Advanced {
		t.REPLAdvanced()
	} else {
		t.REPLBasic()
	}

}
