package main

import (
	"bytes"
	"strings"
	"testing"

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
		"group 1: 2d6+3 Modifier=+3 results=[1 1] total=5",
		"group 2: 1d8 results=[1] total=1",
		"randomness source: test pseudo-random source",
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
		"1d20adv+8 Modifier=+8 results=[4 17] kept=17 total=25",
		"1d20dis-1 Modifier=-1 results=[18 5] kept=5 total=4",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("output %q does not contain %q", output, expected)
		}
	}
}
