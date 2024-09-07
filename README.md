# tierlist
silly repl terminal tierlist creator I created out of boredom (and i did not like online one /shrug)<br>
# Guide on how to use this thing
## Cmd options:
### -mc
Number of spaces you want to add on each side of the longest item in the table,<br> i.e if it was `|item|`, then now it is `|<spaces>item<spaces>|`
### -mt
Same as mc but for tiers, i.e `|tiername|` -> `|<spaces>tiername<spaces>|`
### -a
Alignment mode, 0 by default - centre, 1 - left, 2 - right<br>
Note that if you couple this with mc/mt it will add spaces on the opposite side instead of left and right
### -d
Advanced mode, instead of typing, for example
```
i
s 0
i
a 1
```
you'd be typing
```
i
s 0
a 1
...
q
```
this works with every command that needs additional user input
### -p
Tiers presets, does not work right now, will add later; presets (ignore for now):<br>
<table>
    <tr>
        <th>Number</th>
        <th>Desc.</th>
        <th>Tier names</th>
    </tr>
    <tr>
        <th>0</th>
        <th>Default preset, no tiers</th>
        <th>-</th>
    </tr>
    <tr>
        <th>1</th>
        <th>Basic tiers</th>
        <th>s, a, b, c, d</th>
    </tr>
    <tr>
        <th>2</th>
        <th>Slightly advanced</th>
        <th>s, a, b, c, d, e, f</th>
    </tr>
    <tr>
        <th>3</th>
        <th>If you need really many tiers</th>
        <th>sss, ss, s, a, b, c, d, e, f</th>
    </tr>
</table>

## IO loop aka what can you type and what does that do:
### i
Insert tier
```
i
<tier name> <tier priority, lower is better>
```
### a
Add item to the tier
```
a
<what to add> <tier priority>
```
### rt
Removes tier by either name or priority
```
rt
<tier name> <tier priority>
```

### ri
Removes item from tier if it finds it
```
ri
<item name> <priority>
```

### s
Shows tierlist table
### q
Quits.

# TODO
<ul>
<li>default tiers</li>
</ul>
