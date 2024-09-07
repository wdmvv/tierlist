package main

import (
	"fmt"
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
	iarRegex  *regexp.Regexp = regexp.MustCompile(`(?P<name>[\w\s]+)\s(?P<priority>[0-9]+)`)
	cmdsInput *regexp.Regexp = regexp.MustCompile(`(^[iIaArR]{1}$)|(^(?i)((ri)|(rt))(?-i)$)`)
	cmdsOther *regexp.Regexp = regexp.MustCompile(`^[sSqQ]{1}$`)
)

// basic repl mode that processes one command at a time
func (t *TierList) REPLBasic() {
	for {
		cmdinput := make([]byte, 4)

		_, err := t.input.Read(cmdinput)
		if err != nil {
			if err != io.EOF {
				t.LogErr(err)
			}
			continue
		}
		// trimming \n's and detecting last letter so it does actually match
		ls := 0
		for i, j := range cmdinput {
			if j == 10 {
				cmdinput[i] = 0
			}
			if j != 0 {
				ls = i
			}
		}
		cmdinput = cmdinput[:ls]

		// the horrors of processing commands (such is nature of loops :pensive:)
		if cmdsInput.Match(cmdinput) {
			res := strings.ToLower(string(cmdsInput.Find(cmdinput)))

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
			name := strings.Trim(match[indname], " \n")
			p := match[indp]
			pint, err := strconv.Atoi(p)
			if err != nil {
				t.LogErr(err)
				continue
			}

			switchCmd(t, res, name, pint)
		} else if cmdsOther.Match(cmdinput) {
			res := string(cmdsOther.Find(cmdinput))
			res = strings.ToLower(res)
			switch res {
			case "s":
				t.Show()
			case "q":
				return
			}
		} else {
			fmt.Printf("failed to match command %s, try again\n", string(cmdinput))
		}
	}
}

// advanced mode with loops inside of loop for faster inputs
func (t *TierList) REPLAdvanced() {
	// main repl loop
	for {
		// cuz right now it will never exceed 4 (even 2, but still)
		cmdinput := make([]byte, 4)

		_, err := t.input.Read(cmdinput)
		if err != nil {
			if err != io.EOF {
				t.LogErr(err)
			}
			continue
		}

		ls := 0
		for i, j := range cmdinput {
			if j == 10 {
				cmdinput[i] = 0
			}
			if j != 0 {
				ls = i
			}
		}
		cmdinput = cmdinput[:ls]

		if cmdsInput.Match(cmdinput) {
			res := strings.ToLower(string(cmdsInput.Find(cmdinput)))
			// advanced subloop
			for {
				argsinp := make([]byte, 1024)
				_, err := t.input.Read(argsinp)
				if err != nil {
					if err != io.EOF {
						t.LogErr(err)
					}
					continue
				}
				// trimming input
				ls := 0
				for i, j := range argsinp {
					if j != 0 {
						ls = i
					}
				}
				argsinp = argsinp[:ls]
				// breaking if input is q
				// i really should not check it like this but h
				argsstr := strings.TrimSpace(string(argsinp))

				if (argsinp[0] == 113 || argsinp[0] == 81) && len(argsstr) < 2 {
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
				switchCmd(t, res, name, pint)
			}
		} else if cmdsOther.Match(cmdinput) {
			res := string(cmdsOther.Find(cmdinput))
			res = strings.ToLower(res)
			switch res {
			case "s":
				t.Show()
			case "q":
				return
			}
		} else {
			fmt.Printf("failed to match command %s, try again\n", string(cmdinput))
		}
	}
}

// small func to not repeat same switch twice (plus ill have to edit only once in the future)
func switchCmd(t *TierList, cmd string, name string, priority int) {
	switch cmd {
	case "i":
		t.InsertTier(name, priority)
	case "a":
		t.Add(name, priority)
	case "ri":
		t.RemoveItem(name, priority)
	case "rt":
		t.RemoveTier(name, priority)
	}
}
