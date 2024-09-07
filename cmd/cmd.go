package cmd

import "flag"

type stgs struct {
	MarginItems *int
	MarginTier  *int
	Align       *int
	Advanced    *bool
	Preset      *int
}

// if this was lower(deeper) in the project then i couldve accessed it

var Settings stgs

func ArgsParse() {
	Settings = stgs{
		flag.Int("mc", 0, "how many spaces you want to add on each side of the longest item in tiers"),
		flag.Int("mt", 0, "how many spaces you want to add on each side of the longest tier name"),
		flag.Int("a", 0, "alignment mode for items, 0 - centre, 1 - left 2 - right"),
		flag.Bool("d", false, "advanced mode, i/a/r open subloop instead of jumping to the next command"),
		flag.Int("p", 0, "tiers preset, by default no tiers are created, refer to documentation for tier presets"),
	}
	flag.Parse()
}
