package dice

import (
	"errors"
	"fmt"

	"so-random/internal/randomness"
)

var ErrInvalidDieResult = errors.New("randomness source returned an invalid die result")

func (group DiceGroup) String() string {
	return fmt.Sprintf("%dd%d", group.Count, group.Sides)
}

func (group *DiceGroup) Roll(table randomness.Table) error {
	if group.Count <= 0 || !isSupportedDieSize(uint64(group.Sides)) {
		return ErrInvalidExpression
	}
	if table.PseudoRandom.Roll == nil {
		return randomness.ErrSourceUnavailable
	}

	results := make([]int, 0, group.Count)
	total := 0
	for range group.Count {
		result, err := table.PseudoRandom.Roll(group.Sides)
		if err != nil {
			return fmt.Errorf("pseudo-random source: %w", err)
		}
		if result < 1 || result > group.Sides {
			return ErrInvalidDieResult
		}
		results = append(results, result)
		total += result
	}

	group.Results = results
	group.Rolledvalue = total
	group.SourceDescriptions = make([]string, group.Count)
	for index := range group.SourceDescriptions {
		group.SourceDescriptions[index] = table.PseudoRandom.Description
	}
	return nil
}
