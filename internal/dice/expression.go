package dice

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	MaxExpressionGroups = 20
	MaxDicePerGroup     = 100
	MaxDicePerRequest   = 1000
)

var (
	ErrInvalidExpression = errors.New("invalid dice expression")
	ErrTooManyGroups     = errors.New("too many dice groups")
	ErrTooManyDice       = errors.New("dice count exceeds limit")
)

type DiceGroup struct {
	Count              int
	Sides              int
	Results            []int
	Rolledvalue        int
	SourceDescriptions []string
}

func ParseExpression(expression string) ([]DiceGroup, error) {
	input := strings.ToLower(strings.TrimSpace(expression))
	if input == "" {
		return nil, ErrInvalidExpression
	}

	terms, err := splitTerms(input)
	if err != nil {
		return nil, err
	}
	if len(terms) > MaxExpressionGroups {
		return nil, ErrTooManyGroups
	}

	groups := make([]DiceGroup, 0, len(terms))
	totalDice := 0
	for _, term := range terms {
		group, err := parseTerm(term)
		if err != nil {
			return nil, err
		}
		if group.Count > MaxDicePerGroup {
			return nil, ErrTooManyDice
		}
		if totalDice+group.Count > MaxDicePerRequest {
			return nil, ErrTooManyDice
		}

		totalDice += group.Count
		groups = append(groups, group)
	}

	return groups, nil
}

func splitTerms(input string) ([]string, error) {
	terms := make([]string, 0)
	var term strings.Builder
	seenTerm := false
	commaPending := false

	flushTerm := func() {
		if term.Len() == 0 {
			return
		}
		terms = append(terms, term.String())
		term.Reset()
		seenTerm = true
		commaPending = false
	}

	for _, character := range input {
		switch {
		case character == ',':
			flushTerm()
			if !seenTerm || commaPending {
				return nil, ErrInvalidExpression
			}
			commaPending = true
		case unicode.IsSpace(character):
			flushTerm()
		default:
			term.WriteRune(character)
		}
	}
	flushTerm()
	if commaPending || len(terms) == 0 {
		return nil, ErrInvalidExpression
	}

	return terms, nil
}

func parseTerm(term string) (DiceGroup, error) {
	separator := strings.IndexByte(term, 'd')
	if separator < 0 || separator != strings.LastIndexByte(term, 'd') {
		return DiceGroup{}, ErrInvalidExpression
	}

	count := uint64(1)
	if separator > 0 {
		parsedCount, err := parsePositiveDecimal(term[:separator])
		if err != nil {
			return DiceGroup{}, err
		}
		count = parsedCount
	}

	sides, err := parsePositiveDecimal(term[separator+1:])
	if err != nil || !isSupportedDieSize(sides) {
		return DiceGroup{}, ErrInvalidExpression
	}
	if count > MaxDicePerGroup {
		return DiceGroup{}, ErrTooManyDice
	}

	return DiceGroup{Count: int(count), Sides: int(sides)}, nil
}

func parsePositiveDecimal(value string) (uint64, error) {
	if value == "" {
		return 0, ErrInvalidExpression
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return 0, ErrInvalidExpression
		}
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		return 0, ErrInvalidExpression
	}
	return parsed, nil
}

func isSupportedDieSize(sides uint64) bool {
	switch sides {
	case 2, 4, 6, 8, 10, 12, 20, 100:
		return true
	default:
		return false
	}
}
