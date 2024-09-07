package tlist

import (
	"fmt"
	"sort"
	"strings"
	"tierlist/cmd"
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
// input: i
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
// input: a
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
// input: ri
func (t *TierList) RemoveItem(name string, priority int) {
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

// removes a tier by either priority or name, removes first match(!!!)
// input: rt
func (t *TierList) RemoveTier(name string, priority int) {
	for _, i := range t.Tiers {
		if i.Name == name || i.Priority == priority {
			delete(t.Tiers, i.Name)
			return
		}
	}
}

// display the table
// input: s
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
			*cmd.Settings.MarginTier*2+*cmd.Settings.MarginItems*2) + "\n")

		if len(t.Items) == 0 {
			l, r := centre(t.Name, lngname)
			// tier cell
			out += ("|" + strings.Repeat(" ", l+*cmd.Settings.MarginTier) + t.Name +
				strings.Repeat(" ", r+*cmd.Settings.MarginTier) + "|")
			// items cell
			out += (strings.Repeat(" ", lngitem+*cmd.Settings.MarginItems*2) + "|" + "\n")
			continue
		}
		// what row should we center tier name on/in, top position is prioritized
		namec := len(t.Items) / 2

		for i, j := range t.Items {
			// tier cell
			if i != namec {
				out += "|" + strings.Repeat(" ", lngname+2**cmd.Settings.MarginTier) + "|"
			} else {
				l, r := centre(t.Name, lngname)
				out += ("|" + strings.Repeat(" ", l+*cmd.Settings.MarginTier) + t.Name +
					strings.Repeat(" ", r+*cmd.Settings.MarginTier) + "|")
			}
			// items cell
			l, r := centre(j, lngitem)

			// left/centre/right alignment
			switch *cmd.Settings.Align {
			// centre
			case 0:
				out += (strings.Repeat(" ", l+*cmd.Settings.MarginItems) + j +
					strings.Repeat(" ", r+*cmd.Settings.MarginItems) + "|\n")
			// left
			case 1:
				out += (j + strings.Repeat(" ", l+r+*cmd.Settings.MarginItems*2) + "|\n")
			case 2:
				out += (strings.Repeat(" ", l+r+*cmd.Settings.MarginItems*2) + j + "|\n")
			}
		}
	}
	// closing border
	out += (strings.Repeat("-", lngname+lngitem+3+
		*cmd.Settings.MarginItems*2+*cmd.Settings.MarginTier*2) + "\n")

	t.output.Write([]byte(out))
}

// TODO: generates preset
func (t *TierList) GenPreset(n int) error {
	if n == 0 {
		return nil
	}
	var ts []string
	switch n {
	case 1:
		ts = []string{"s", "a", "b", "c", "d"}
	case 2:
		ts = []string{"s", "a", "b", "c", "d", "e", "f"}
	case 3:
		ts = []string{"sss", "ss", "s", "a", "b", "c", "d", "e", "f"}
	default:
		return fmt.Errorf("unknown preset %d, ignoring...", n)
	}
	for i, j := range ts {
		t.InsertTier(j, i)
	}
	return nil
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
