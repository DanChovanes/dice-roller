package main

import (
	"fmt"

	"so-random/internal/dice"
	"so-random/internal/randomness"
)

func main() {
	groups, err := dice.ParseExpression("4d6, d8 d20")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	table := randomness.NewTable()
	for index := range groups {
		group := &groups[index]
		if err := group.Roll(table); err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Printf("group %d: %s results=%v total=%d\n", index+1, group, group.Results, group.Rolledvalue)
	}
}
