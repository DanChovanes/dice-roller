package main

import (
	"bytes"
	"strings"
	"testing"

	"so-random/internal/dice"
	"so-random/internal/randomness"
)

func TestRunAcceptsExpressionArgumentsAndShowsSourceInformation(t *testing.T) {
	table := randomness.Table{
		PseudoRandom: randomness.Source{
			Description: "test pseudo-random source",
			Roll:        func(int) (int, error) { return 1, nil },
		},
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-i", "2d6+3", "d8"}, table, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() = %d, stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"The dice have spoken!",
		"1. 2d6+3",
		"rolls: [1] [1]",
		"modifier: +3",
		"TOTAL: 5",
		"2. 1d8",
		"TOTAL: 1",
		"A whisper from the source:",
		"test pseudo-random source",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("output %q does not contain %q", output, expected)
		}
	}
}

func TestRunRequiresExpression(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(nil, randomness.NewTable(), &stdout, &stderr)
	if code == 0 {
		t.Fatal("run() = 0, want usage error")
	}
	if !strings.Contains(stderr.String(), "Usage:") {
		t.Fatalf("stderr = %q, want usage", stderr.String())
	}
}

func TestRunHelpShowsQuickStartAndNotation(t *testing.T) {
	for _, helpFlag := range []string{"--help", "-h"} {
		t.Run(helpFlag, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := run([]string{helpFlag}, randomness.Table{}, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("run() = %d, want 0; stderr=%q", code, stderr.String())
			}
			for _, expected := range []string{
				"Usage: whim [OPTIONS] DICE...",
				"QUICK START",
				"whim d20",
				"whim 4d6+2",
				"DICE NOTATION",
				"adv / dis",
				"Dice sizes:",
				"--color MODE",
			} {
				if !strings.Contains(stdout.String(), expected) {
					t.Errorf("help output %q does not contain %q", stdout.String(), expected)
				}
			}
			if stderr.Len() != 0 {
				t.Errorf("stderr = %q, want empty", stderr.String())
			}
		})
	}
}

func TestRunDoesNotPrintPartialBatchOnSourceFailure(t *testing.T) {
	callCount := 0
	table := randomness.Table{
		PseudoRandom: randomness.Source{
			Description: "test source",
			Roll: func(int) (int, error) {
				callCount++
				if callCount == 2 {
					return 0, randomness.ErrSourceUnavailable
				}
				return 1, nil
			},
		},
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"d6", "d8"}, table, &stdout, &stderr)
	if code == 0 {
		t.Fatal("run() = 0, want source error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want no partial batch output", stdout.String())
	}
}

func TestRunShowsBothAdvantageRollsKeptValueAndModifier(t *testing.T) {
	results := []int{4, 17, 18, 5}
	resultIndex := 0
	table := randomness.Table{
		PseudoRandom: randomness.Source{
			Description: "test source",
			Roll: func(int) (int, error) {
				result := results[resultIndex]
				resultIndex++
				return result, nil
			},
		},
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"d20+8adv", "d20-1dis"}, table, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run() = %d, stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"1d20adv+8",
		"rolls: [4] (discarded) [17] (kept)",
		"kept: 17 (advantage)",
		"modifier: +8",
		"TOTAL: 25",
		"1d20dis-1",
		"rolls: [18] (discarded) [5] (kept)",
		"kept: 5 (disadvantage)",
		"modifier: -1",
		"TOTAL: 4",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("output %q does not contain %q", output, expected)
		}
	}
}

func TestRunColorModeControlsANSIOutput(t *testing.T) {
	for _, test := range []struct {
		name      string
		mode      string
		wantColor bool
	}{
		{name: "always", mode: "always", wantColor: true},
		{name: "never", mode: "never", wantColor: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := run([]string{"--color=" + test.mode, "d6"}, randomness.NewTable(), &stdout, &stderr)
			if code != 0 {
				t.Fatalf("run() = %d, stderr=%q", code, stderr.String())
			}
			gotColor := strings.Contains(stdout.String(), "\x1b[")
			if gotColor != test.wantColor {
				t.Errorf("ANSI output = %t, want %t; output=%q", gotColor, test.wantColor, stdout.String())
			}
		})
	}
}

func TestColorEnabledPolicy(t *testing.T) {
	for _, test := range []struct {
		name     string
		mode     string
		terminal bool
		noColor  bool
		want     bool
	}{
		{name: "auto terminal", mode: "auto", terminal: true, want: true},
		{name: "auto redirected", mode: "auto", want: false},
		{name: "auto honors NO_COLOR", mode: "auto", terminal: true, noColor: true, want: false},
		{name: "always overrides NO_COLOR", mode: "always", noColor: true, want: true},
		{name: "never", mode: "never", terminal: true, want: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := colorEnabled(test.mode, test.terminal, test.noColor); got != test.want {
				t.Errorf("colorEnabled() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestFormatRollsMarksAdvantageTies(t *testing.T) {
	group := dice.DiceGroup{
		Advantage:  dice.AdvantageOn,
		KeptResult: 17,
		Results:    []int{17, 17},
	}
	if got, want := formatRolls(group, false), "[17] (kept) [17] (tied)"; got != want {
		t.Errorf("formatRolls() = %q, want %q", got, want)
	}
}

func TestRunCanColorErrorsWhenForced(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	table := randomness.Table{}
	code := run([]string{"--color=always", "d6"}, table, &stdout, &stderr)
	if code == 0 {
		t.Fatal("run() = 0, want source error")
	}
	if !strings.Contains(stderr.String(), "\x1b[1;31merror:") {
		t.Fatalf("stderr = %q, want colored error prefix", stderr.String())
	}
}

func TestRunRejectsUnknownColorMode(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := run([]string{"--color=bright", "d6"}, randomness.NewTable(), &stdout, &stderr)
	if code == 0 {
		t.Fatal("run() = 0, want invalid color mode error")
	}
	if !strings.Contains(stderr.String(), "invalid color mode") {
		t.Fatalf("stderr = %q, want invalid color mode error", stderr.String())
	}
}
