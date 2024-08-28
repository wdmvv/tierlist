package main

import (
	"sort"
	"strings"
)

func (t *TierList) LogErr(e error) {
	t.output.Write([]byte(e.Error()))
}

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
	t.Tiers[name] = &Tier{priority, name, make([]string, 0)}
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
	tiers := make([]Tier, 0, len(t.Tiers))

	for _, i := range t.Tiers {
		tiers = append(tiers, *i)
	}
	sortTiers(tiers)

	var (
		lngname int
		lngitem int
	)

	for _, i := range tiers {
		if len(i.Name) > lngname {
			lngname = len(i.Name)
		}
		for _, j := range i.Items {

			if len(j) > lngitem {
				lngitem = len(j)
			}
		}
	}

	// i wonder if there is a better way
	// maybe i should show it as i create it,
	// maybe i should display on row basis
	// maybe buffering of sorts?
	// who knows what is the best way

	var out string
	for _, t := range tiers {
		out += (strings.Repeat("-", lngname+lngitem+3+
			*Settings.MarginTier*2+*Settings.MarginItems*2) + "\n")

		if len(t.Items) == 0 {
			l, r := centre(t.Name, lngname)
			// tier cell
			out += ("|" + strings.Repeat(" ", l+*Settings.MarginTier) + t.Name +
				strings.Repeat(" ", r+*Settings.MarginTier) + "|")
			// items cell
			out += (strings.Repeat(" ", lngitem+*Settings.MarginItems*2) + "|" + "\n")
			continue
		}
		// what row should we center tier name on/in, top position is prioritized
		namec := len(t.Items) / 2

		for i, j := range t.Items {
			// tier cell
			if i != namec {
				out += "|" + strings.Repeat(" ", lngname+2**Settings.MarginTier) + "|"
			} else {
				l, r := centre(t.Name, lngname)
				out += ("|" + strings.Repeat(" ", l+*Settings.MarginTier) + t.Name +
					strings.Repeat(" ", r+*Settings.MarginTier) + "|")
			}
			// items cell
			l, r := centre(j, lngitem)

			// l/c/r alignment
			switch *Settings.Align {
			// centre
			case 0:
				out += (strings.Repeat(" ", l+*Settings.MarginItems) + j +
					strings.Repeat(" ", r+*Settings.MarginItems) + "|\n")
			// left
			case 1:
				out += (j + strings.Repeat(" ", l+r+*Settings.MarginItems*2) + "|\n")
			case 2:
				out += (strings.Repeat(" ", l+r+*Settings.MarginItems*2) + j + "|\n")
			}
		}
	}
	// closing border
	out += (strings.Repeat("-", lngname+lngitem+3+
		*Settings.MarginItems*2+*Settings.MarginTier*2) + "\n")

	t.output.Write([]byte(out))
}

// centre text based on the box size
// smth is the string to be centered, lng is the max width
func centre(smth string, lng int) (int, int) {
	var left, right int
	if len(smth) == lng {
		return 0, 0
	}
	// all avlb. distance by 2
	left = (lng - len(smth)) / 2
	// leftovers (or rightovers ^^ )
	// _presumably_ this is bigger that 0
	right = lng - len(smth) - left
	return left, right
}

func sortTiers(ts []Tier) {
	sort.SliceStable(ts, func(i, j int) bool {
		return ts[i].Priority < ts[j].Priority
	})
}
