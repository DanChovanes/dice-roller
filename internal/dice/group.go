package dice

import (
	"errors"
	"fmt"
	"strconv"

	"so-random/internal/randomness"
)

var ErrInvalidDieResult = errors.New("randomness source returned an invalid die result")

func (group DiceGroup) String() string {
	expression := fmt.Sprintf("%dd%d", group.Count, group.Sides)
	if group.Advantage == AdvantageOn {
		expression += "adv"
	} else if group.Advantage == DisadvantageOn {
		expression += "dis"
	}
	return expression + group.Modifier
}

func (group *DiceGroup) Roll(table randomness.Table) error {
	if group.Count <= 0 || group.Count > MaxDicePerGroup || !isSupportedDieSize(uint64(group.Sides)) {
		return ErrInvalidExpression
	}
	if group.Advantage != NoAdvantage && (group.Sides != 20 || group.Count != 1) {
		return ErrInvalidExpression
	}
	if table.PseudoRandom.Roll == nil {
		return randomness.ErrSourceUnavailable
	}

	modifier := 0
	if group.Modifier != "" {
		parsedModifier, err := strconv.Atoi(group.Modifier)
		if err != nil {
			return ErrInvalidExpression
		}
		modifier = parsedModifier
	}

	rollCount := group.Count
	if group.Advantage != NoAdvantage {
		rollCount = 2
	}
	results := make([]int, 0, rollCount)
	sources := make([]string, 0, rollCount)
	total := 0
	keptResult := 0
	for index := range rollCount {
		result, err := table.PseudoRandom.Roll(group.Sides)
		if err != nil {
			return fmt.Errorf("pseudo-random source: %w", err)
		}
		if result < 1 || result > group.Sides {
			return ErrInvalidDieResult
		}
		results = append(results, result)
		sources = append(sources, table.PseudoRandom.Description)
		if index == 0 || (group.Advantage == AdvantageOn && result > keptResult) ||
			(group.Advantage == DisadvantageOn && result < keptResult) {
			keptResult = result
		}
		if group.Advantage == NoAdvantage {
			total += result
		}
	}

	if group.Advantage != NoAdvantage {
		total = keptResult
		group.KeptResult = keptResult
	} else {
		group.KeptResult = 0
	}
	group.Results = results
	group.Rolledvalue = total + modifier
	group.SourceDescriptions = sources
	return nil
}
