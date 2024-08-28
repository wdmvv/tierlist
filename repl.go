package main

import (
	"io"
	"regexp"
	"strconv"
	"strings"
)

// my thoughts on how this will(should) work:
// you can either give own tiers or use default i.e D-S (or smth)
// additionally you can insert own tier as well but you have to specify position where 0 is the highest N is the lowest
// TierList is using string:tier for somewhat? faster lookup (although size is small so one could assume list is o(1) too)
// so it is for more convenient lookup i guess (i would've went with list where index is a priority)
// will just stick with this for now, might change later
//
// somehow regex should support either tier height or tier name, i really should think about this
// and i should normalize capitals i.e S s are the same, so maybe .lower() check?
// i guess ill have to find out more about groups

type TierList struct {
	Tiers  map[string]*Tier
	input  io.Reader
	output io.Writer
	lowest int // lowest tier index
}

func NewTierlist(i io.Reader, o io.Writer) *TierList {
	return &TierList{make(map[string]*Tier), i, o, 0}
}

type Tier struct {
	Priority int
	Name     string
	Items    []string
}

var (
	iarRegex *regexp.Regexp = regexp.MustCompile(`(?P<name>[\w\s]+)\s(?P<priority>[0-9]+)`)
)

// basic repl mode that processes one command at a time
func (t *TierList) REPLBasic() {
	// ill leave this as priorities and not tier names because i cant come up with a good idea of how to manage this in single line input (so its 2 now)
	// insert - ^([iI]{1})\s(?P<item>[\w\s]+)\s(?P<priority>[0-9]+) - <command> <tier_name> <int priority>
	// add - ^([aA]{1})\s(?P<item>[\w\s]+)\s(?P<priority>[0-9]+) - <command> <item> <int priority>
	// remove - ^([rR]{1})\s(?P<item>[\w\s]+)\s(?P<priority>[0-9]+) - <command> <item> <int priority>
	// show - ^([sS]{1}) - <command> (maybe will add more params later)
	// quit - ^([qQ]{1}) - <command>

	cmds := regexp.MustCompile(`^[iIaArRsSqQ]{1}`)

	for {
		cmdinput := make([]byte, 1024)

		_, err := t.input.Read(cmdinput)
		if err != nil {
			if err != io.EOF {
				t.LogErr(err)
			}
			continue
		}
		res := cmds.Find(cmdinput)
		if res == nil {
			continue
		}

		// the horrors of processing commands (sorry!)
		//
		switch string(res) {
		case "i", "I", "a", "A", "r", "R":
			argsinp := make([]byte, 1024)
			_, err = t.input.Read(argsinp)
			if err != nil {
				if err != io.EOF {
					t.LogErr(err)
				}
				continue
			}

			line := string(argsinp)
			match := iarRegex.FindStringSubmatch(line)
			indname := iarRegex.SubexpIndex("name")
			indp := iarRegex.SubexpIndex("priority")
			if indname >= len(match) || indp >= len(match) || indname == -1 || indp == -1 {
				continue
			}
			name := match[indname]
			p := match[indp]
			pint, err := strconv.Atoi(p)
			if err != nil {
				t.LogErr(err)
				continue
			}
			switch string(res) {
			case "i", "I":
				t.InsertTier(name, pint)
			case "a", "A":
				t.Add(name, pint)
			case "r", "R":
				t.Remove(name, pint)
			}
		case "s", "S":
			t.Show()
		case "q", "Q":
			return
		}
	}
}

// advanced mode with loops inside of loop for faster inputs
func (t *TierList) REPLAdvanced() {
	cmds := regexp.MustCompile(`^[iIaArRsSqQ]{1}`)

	// main repl loop
	for {
		cmdinput := make([]byte, 1024)

		_, err := t.input.Read(cmdinput)
		if err != nil {
			if err != io.EOF {
				t.LogErr(err)
			}
			continue
		}
		res := cmds.Find(cmdinput)
		if res == nil {
			continue
		}

		switch string(res) {
		case "i", "I", "a", "A", "r", "R":
			// subloop for these commands
			for {
				argsinp := make([]byte, 1024)
				_, err = t.input.Read(argsinp)
				if err != nil {
					if err != io.EOF {
						t.LogErr(err)
					}
					continue
				}

				argsstr := strings.TrimSpace(string(argsinp))

				// i tried to do this with regex and failed, thus this stupid workaround
				if argsstr[0] == 113 || argsstr[0] == 81 {
					break
				}

				line := string(argsinp)
				match := iarRegex.FindStringSubmatch(line)
				indname := iarRegex.SubexpIndex("name")
				indp := iarRegex.SubexpIndex("priority")
				if indname >= len(match) || indp >= len(match) || indname == -1 || indp == -1 {
					continue
				}
				name := match[indname]
				p := match[indp]
				pint, err := strconv.Atoi(p)
				if err != nil {
					t.LogErr(err)
					continue
				}
				switch string(res) {
				case "i", "I":
					t.InsertTier(name, pint)
				case "a", "A":
					t.Add(name, pint)
				case "r", "R":
					t.Remove(name, pint)
				}
			}
		case "s", "S":
			t.Show()
		case "q", "Q":
			return
		}
	}
}
