package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"so-random/internal/dice"
	"so-random/internal/randomness"
)

func main() {
	os.Exit(run(os.Args[1:], randomness.NewTable(), os.Stdout, os.Stderr))
}

func run(args []string, table randomness.Table, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("whim", flag.ContinueOnError)
	flags.SetOutput(stderr)
	showInformation := flags.Bool("i", false, "show descriptions of randomness sources used")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: whim [-i] <dice-expression...>")
		fmt.Fprintln(stderr, "Example: whim -i 4d6 d8")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 2
	}

	expression := strings.Join(flags.Args(), " ")
	if expression == "" {
		flags.Usage()
		return 2
	}

	groups, err := dice.ParseExpression(expression)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 2
	}

	for index := range groups {
		group := &groups[index]
		if err := group.Roll(table); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}

	for index := range groups {
		group := &groups[index]
		modifierOutput := ""
		if group.Modifier != "" {
			modifierOutput = fmt.Sprintf(" Modifier=%s", group.Modifier)
		}
		if group.Advantage == dice.NoAdvantage {
			fmt.Fprintf(stdout, "group %d: %s%s results=%v total=%d\n", index+1, group, modifierOutput, group.Results, group.Rolledvalue)
		} else {
			fmt.Fprintf(stdout, "group %d: %s%s results=%v kept=%d total=%d\n", index+1, group, modifierOutput, group.Results, group.KeptResult, group.Rolledvalue)
		}
	}
	if *showInformation {
		for _, description := range usedSourceDescriptions(groups) {
			fmt.Fprintf(stdout, "randomness source: %s\n", description)
		}
	}
	return 0
}

func usedSourceDescriptions(groups []dice.DiceGroup) []string {
	seen := make(map[string]struct{})
	descriptions := make([]string, 0)
	for _, group := range groups {
		for _, description := range group.SourceDescriptions {
			if description == "" {
				description = "description unavailable"
			}
			if _, exists := seen[description]; exists {
				continue
			}
			seen[description] = struct{}{}
			descriptions = append(descriptions, description)
		}
	}
	return descriptions
}
