package main

import (
	"os"
	"tierlist/cmd"
	"tierlist/pkg/tlist"
)

func main() {

	cmd.ArgsParse()
	t := tlist.NewTierlist(os.Stdin, os.Stdout)
	err := t.GenPreset(*cmd.Settings.Preset)
	if err != nil {
		t.LogErr(err)
	}
	if *cmd.Settings.Advanced {
		t.REPLAdvanced()
	} else {
		t.REPLBasic()
	}

}
