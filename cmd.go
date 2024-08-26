package main

import "flag"

type stgs struct {
	MarginCell *int
	MarginTier *int
	Align      *int
	Advanced   *bool
}

// if this was lower(deeper) in the project then i couldve accessed it

var Settings stgs

func ArgsParse() {
	Settings = stgs{
		flag.Int("mc", 0, "how many spaces you want it to have on both sides of the cell's longest item"),
		flag.Int("mt", 0, "how many spaces you want it to have on both sides of the tier's longest item"),
		flag.Int("a", 0, "alignment mode, 0 - centre, 1 - left 2 - right"),
		flag.Bool("d", false, "advanced mode, i/a/r open subloop instead of jumping to the next command"),
	}
	flag.Parse()
}
