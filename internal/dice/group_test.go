package dice

import (
	"errors"
	"reflect"
	"testing"

	"so-random/internal/randomness"
)

func TestDiceGroupRollStoresResultsAndTotal(t *testing.T) {
	group := DiceGroup{Count: 4, Sides: 6}
	results := []int{2, 5, 1, 4}
	resultIndex := 0
	table := randomness.Table{
		PseudoRandom: randomness.Source{
			Description: "deterministic test source",
			Roll: func(sides int) (int, error) {
				if sides != 6 {
					t.Fatalf("source received sides = %d, want 6", sides)
				}
				result := results[resultIndex]
				resultIndex++
				return result, nil
			},
		},
	}

	if err := group.Roll(table); err != nil {
		t.Fatalf("Roll() error = %v", err)
	}
	if !reflect.DeepEqual(group.Results, results) {
		t.Fatalf("Results = %v, want %v", group.Results, results)
	}
	if group.Rolledvalue != 12 {
		t.Fatalf("Rolledvalue = %d, want 12", group.Rolledvalue)
	}
	if !reflect.DeepEqual(group.SourceDescriptions, []string{"deterministic test source", "deterministic test source", "deterministic test source", "deterministic test source"}) {
		t.Fatalf("SourceDescriptions = %v", group.SourceDescriptions)
	}
}

func TestDiceGroupRollDoesNotStorePartialResultsOnSourceFailure(t *testing.T) {
	group := DiceGroup{Count: 2, Sides: 6, Results: []int{3}, Rolledvalue: 3}
	callCount := 0
	table := randomness.Table{
		PseudoRandom: randomness.Source{
			Description: "failing test source",
			Roll: func(int) (int, error) {
				callCount++
				if callCount == 2 {
					return 0, errors.New("source failed")
				}
				return 4, nil
			},
		},
	}

	if err := group.Roll(table); err == nil {
		t.Fatal("Roll() error = nil, want source failure")
	}
	if !reflect.DeepEqual(group.Results, []int{3}) || group.Rolledvalue != 3 {
		t.Fatalf("failed roll modified group: results=%v total=%d", group.Results, group.Rolledvalue)
	}
}

func TestDiceGroupRollRejectsInvalidSourceResult(t *testing.T) {
	group := DiceGroup{Count: 1, Sides: 6}
	table := randomness.Table{
		PseudoRandom: randomness.Source{
			Description: "invalid test source",
			Roll:        func(int) (int, error) { return 7, nil },
		},
	}

	if err := group.Roll(table); !errors.Is(err, ErrInvalidDieResult) {
		t.Fatalf("Roll() error = %v, want %v", err, ErrInvalidDieResult)
	}
}

func TestDiceGroupRollAddsModifierToAllDiceTotal(t *testing.T) {
	group := DiceGroup{Count: 2, Sides: 4, Modifier: "-3"}
	results := []int{4, 2}
	resultIndex := 0
	table := randomness.Table{PseudoRandom: randomness.Source{
		Description: "deterministic test source",
		Roll: func(int) (int, error) {
			result := results[resultIndex]
			resultIndex++
			return result, nil
		},
	}}

	if err := group.Roll(table); err != nil {
		t.Fatalf("Roll() error = %v", err)
	}
	if group.Rolledvalue != 3 {
		t.Fatalf("Rolledvalue = %d, want 3", group.Rolledvalue)
	}
}

func TestDiceGroupRollAdvantageKeepsHigherD20(t *testing.T) {
	group := DiceGroup{Count: 1, Sides: 20, Modifier: "+8", Advantage: AdvantageOn}
	results := []int{4, 17}
	resultIndex := 0
	table := randomness.Table{PseudoRandom: randomness.Source{
		Description: "deterministic test source",
		Roll: func(sides int) (int, error) {
			if sides != 20 {
				t.Fatalf("source received sides = %d, want 20", sides)
			}
			result := results[resultIndex]
			resultIndex++
			return result, nil
		},
	}}

	if err := group.Roll(table); err != nil {
		t.Fatalf("Roll() error = %v", err)
	}
	if !reflect.DeepEqual(group.Results, results) {
		t.Fatalf("Results = %v, want both rolls %v", group.Results, results)
	}
	if group.KeptResult != 17 || group.Rolledvalue != 25 {
		t.Fatalf("KeptResult=%d Rolledvalue=%d, want kept 17 total 25", group.KeptResult, group.Rolledvalue)
	}
}

func TestDiceGroupRollDisadvantageKeepsLowerD20(t *testing.T) {
	group := DiceGroup{Count: 1, Sides: 20, Modifier: "-1", Advantage: DisadvantageOn}
	results := []int{18, 5}
	resultIndex := 0
	table := randomness.Table{PseudoRandom: randomness.Source{
		Description: "deterministic test source",
		Roll: func(int) (int, error) {
			result := results[resultIndex]
			resultIndex++
			return result, nil
		},
	}}

	if err := group.Roll(table); err != nil {
		t.Fatalf("Roll() error = %v", err)
	}
	if !reflect.DeepEqual(group.Results, results) {
		t.Fatalf("Results = %v, want both rolls %v", group.Results, results)
	}
	if group.KeptResult != 5 || group.Rolledvalue != 4 {
		t.Fatalf("KeptResult=%d Rolledvalue=%d, want kept 5 total 4", group.KeptResult, group.Rolledvalue)
	}
}
