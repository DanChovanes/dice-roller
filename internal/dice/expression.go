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
	Modifier           string
	Advantage          Advantage
	Results            []int
	KeptResult         int
	Rolledvalue        int
	Sides              int
	SourceDescriptions []string
}

type Advantage string

const (
	NoAdvantage    Advantage = ""
	AdvantageOn    Advantage = "advantage"
	DisadvantageOn Advantage = "disadvantage"
)

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
	if separator < 0 {
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

	remainder := term[separator+1:]
	sidesEnd := 0
	for sidesEnd < len(remainder) && remainder[sidesEnd] >= '0' && remainder[sidesEnd] <= '9' {
		sidesEnd++
	}
	sides, err := parsePositiveDecimal(remainder[:sidesEnd])
	if err != nil || !isSupportedDieSize(sides) {
		return DiceGroup{}, ErrInvalidExpression
	}
	if count > MaxDicePerGroup {
		return DiceGroup{}, ErrTooManyDice
	}

	group := DiceGroup{Count: int(count), Sides: int(sides)}
	suffix := remainder[sidesEnd:]
	if strings.HasPrefix(suffix, "adv") || strings.HasPrefix(suffix, "dis") {
		group.Advantage, suffix = parseAdvantage(suffix)
	}
	if strings.HasPrefix(suffix, "+") || strings.HasPrefix(suffix, "-") {
		modifierEnd := 1
		for modifierEnd < len(suffix) && suffix[modifierEnd] >= '0' && suffix[modifierEnd] <= '9' {
			modifierEnd++
		}
		modifier, err := parseModifier(suffix[:modifierEnd])
		if err != nil {
			return DiceGroup{}, err
		}
		group.Modifier = modifier
		suffix = suffix[modifierEnd:]
		if group.Advantage == NoAdvantage {
			group.Advantage, suffix = parseAdvantage(suffix)
		}
	}
	if suffix != "" {
		return DiceGroup{}, ErrInvalidExpression
	}

	if group.Advantage != NoAdvantage && (group.Sides != 20 || group.Count != 1) {
		return DiceGroup{}, ErrInvalidExpression
	}
	return group, nil
}

func parseModifier(value string) (string, error) {
	if len(value) < 2 || (value[0] != '+' && value[0] != '-') {
		return "", ErrInvalidExpression
	}
	for _, character := range value[1:] {
		if character < '0' || character > '9' {
			return "", ErrInvalidExpression
		}
	}

	modifier, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return "", ErrInvalidExpression
	}
	if modifier > 0 {
		return "+" + strconv.FormatInt(modifier, 10), nil
	}
	return strconv.FormatInt(modifier, 10), nil
}

func parseAdvantage(value string) (Advantage, string) {
	if strings.HasPrefix(value, "adv") {
		return AdvantageOn, value[len("adv"):]
	}
	if strings.HasPrefix(value, "dis") {
		return DisadvantageOn, value[len("dis"):]
	}
	return NoAdvantage, value
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
