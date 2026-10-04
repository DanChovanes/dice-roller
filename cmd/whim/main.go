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
	colorMode := flags.String("color", "auto", "color output: auto, always, or never")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: whim [-i] [--color=auto|always|never] <dice-expression...>")
		fmt.Fprintln(stderr, "Example: whim -i 4d6 d8")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 2
	}
	if *colorMode != "auto" && *colorMode != "always" && *colorMode != "never" {
		fmt.Fprintf(stderr, "error: invalid color mode %q (want auto, always, or never)\n", *colorMode)
		return 2
	}
	errorColor := colorEnabled(*colorMode, isTerminal(stderr), os.Getenv("NO_COLOR") != "")

	expression := strings.Join(flags.Args(), " ")
	if expression == "" {
		flags.Usage()
		return 2
	}

	groups, err := dice.ParseExpression(expression)
	if err != nil {
		writeError(stderr, err, errorColor)
		return 2
	}

	for index := range groups {
		group := &groups[index]
		if err := group.Roll(table); err != nil {
			writeError(stderr, err, errorColor)
			return 1
		}
	}

	useColor := colorEnabled(*colorMode, isTerminal(stdout), os.Getenv("NO_COLOR") != "")
	fmt.Fprintln(stdout, paint("The dice have spoken!", "1;36", useColor))
	for index := range groups {
		printGroup(stdout, index+1, groups[index], useColor)
	}
	if *showInformation {
		descriptions := usedSourceDescriptions(groups)
		if len(descriptions) > 0 {
			fmt.Fprintln(stdout)
			fmt.Fprintln(stdout, paint("A whisper from the source:", "2", useColor))
		}
		for _, description := range descriptions {
			fmt.Fprintf(stdout, "  - %s\n", description)
		}
	}
	return 0
}

func printGroup(output io.Writer, number int, group dice.DiceGroup, useColor bool) {
	fmt.Fprintf(output, "\n  %d. %s\n", number, paint(group.String(), "1;36", useColor))
	fmt.Fprintf(output, "     rolls: %s\n", formatRolls(group, useColor))
	if group.Advantage != dice.NoAdvantage {
		advantageLabel := "advantage"
		if group.Advantage == dice.DisadvantageOn {
			advantageLabel = "disadvantage"
		}
		fmt.Fprintf(output, "     kept: %s (%s)\n", paint(fmt.Sprint(group.KeptResult), "1;32", useColor), advantageLabel)
	}
	if group.Modifier != "" {
		fmt.Fprintf(output, "     modifier: %s\n", paint(group.Modifier, "2", useColor))
	}
	fmt.Fprintf(output, "     ==> TOTAL: %s\n", paint(fmt.Sprint(group.Rolledvalue), "1;33", useColor))
}

func formatRolls(group dice.DiceGroup, useColor bool) string {
	rolls := make([]string, len(group.Results))
	keptIndex := -1
	if group.Advantage != dice.NoAdvantage {
		for index, result := range group.Results {
			if result == group.KeptResult {
				keptIndex = index
				break
			}
		}
	}
	for index, result := range group.Results {
		roll := fmt.Sprintf("[%d]", result)
		if keptIndex >= 0 {
			switch {
			case index == keptIndex:
				roll += " (kept)"
			case result == group.KeptResult:
				roll += " (tied)"
			default:
				roll += " (discarded)"
				rolls[index] = paint(roll, "2", useColor)
				continue
			}
		}
		rolls[index] = paint(roll, "1;32", useColor)
	}
	return strings.Join(rolls, " ")
}

func colorEnabled(mode string, terminal, noColor bool) bool {
	switch mode {
	case "always":
		return true
	case "never":
		return false
	default:
		return terminal && !noColor
	}
}

func isTerminal(output io.Writer) bool {
	file, ok := output.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func paint(value, code string, enabled bool) string {
	if !enabled {
		return value
	}
	return "\x1b[" + code + "m" + value + "\x1b[0m"
}

func writeError(output io.Writer, err error, useColor bool) {
	fmt.Fprintf(output, "%s %v\n", paint("error:", "1;31", useColor), err)
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
