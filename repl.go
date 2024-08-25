package main

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
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
	Tiers  map[string]Tier
	input  io.Reader
	output io.Writer
	lowest int // lowest tier index
}

func NewTierlist(i io.Reader, o io.Writer) *TierList {
	return &TierList{make(map[string]Tier), i, o, 0}
}

type Tier struct {
	Priority int
	Name     string
	Items    []string
}

var (
	iarRegex *regexp.Regexp = regexp.MustCompile(`(?P<name>[\w\s]+)\s(?P<priority>[0-9]+)`)
)

func (t *TierList) REPL() {
	// ill leave this as priorities and not tier names because i cant come up with a good idea of how to manage this in single line input (so its 2 now)
	// insert - ^([iI]{1})\s(?P<item>[\w\s]+)\s(?P<priority>[0-9]+) - <command> <tier_name> <int priority>
	// add - ^([aA]{1})\s(?P<item>[\w\s]+)\s(?P<priority>[0-9]+) - <command> <item> <int priority>
	// remove - ^([rR]{1})\s(?P<item>[\w\s]+)\s(?P<priority>[0-9]+) - <command> <item> <int priority>
	// show - ^([sS]{1}) - <command> (maybe will add more params later)
	// quit - ^([qQ]{1}) - <command>

	c, err := regexp.Compile(`^[iI aA rR sS qQ]{1}`)
	if err != nil {
		t.LogErr(err)
	}
	inp := make([]byte, 1024)
	for {
		_, err = t.input.Read(inp)
		if err != nil {
			t.LogErr(err)
			continue
		}
		res := c.Find(inp)
		if res == nil {
			continue
		}

		// the horrors of processing commands (sorry!)

		switch string(res) {
		case "i", "I", "a", "A", "r", "R":
			_, err = t.input.Read(inp)
			if err != nil {
				t.LogErr(err)
				continue
			}
			line := string(inp)
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
			break
		}
	}
}

func (t *TierList) LogErr(e error) {
	t.output.Write([]byte(e.Error()))
}

// imagine generics here... ↑↓

func (t *TierList) LogStr(s string) {
	t.output.Write([]byte(s))
}

// insert new tier, 0 is the highest, N is the lowest
// if priority exists then it updates the name
// if name exists then it ignores (should add swap function maybe)
func (t *TierList) InsertTier(name string, priority int) {
	if _, ok := t.Tiers[name]; ok {
		t.LogStr(name + " tier does not exist!")
		return
	}
	for _, i := range t.Tiers {
		if i.Priority == priority {
			i.Name = name
			return
		}
	}
	t.Tiers[name] = Tier{priority, name, make([]string, 0)}
}

// add item to a tier
func (t *TierList) Add(name string, priority int) {
	for _, i := range t.Tiers {
		if i.Priority == priority {
			i.Items = append(i.Items, name)
			return
		}
	}
	t.LogStr("no such priority found")
}

// remove item from a tier
func (t *TierList) Remove(name string, priority int) {
	for _, i := range t.Tiers {
		if i.Priority == priority {
			for j, k := range i.Items {
				if k == name {
					i.Items = append(i.Items[:j], i.Items[j+1:]...)
					return
				}
			}
		}
	}
	t.LogStr("no item or priority found")
}

// display
func (t *TierList) Show() {
	ts := make([]Tier, len(t.Tiers))

	for _, i := range t.Tiers {
		ts = append(ts, i)
	}
	sortTiers(ts)

	for _, i := range ts {
		fmt.Println(i.Name, i.Items)
	}

}

func sortTiers(ts []Tier) {
	sort.SliceStable(ts, func(i, j int) bool {
		return ts[i].Priority > ts[j].Priority
	})
}
