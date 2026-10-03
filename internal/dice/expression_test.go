package dice

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestParseExpressionValid(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       []DiceGroup
	}{
		{
			name:       "single group",
			expression: "4d6",
			want:       []DiceGroup{{Count: 4, Sides: 6}},
		},
		{
			name:       "omitted count and normalized case",
			expression: " D8 ",
			want:       []DiceGroup{{Count: 1, Sides: 8}},
		},
		{
			name:       "comma separated groups",
			expression: "4d6, 2d10, 1d6, d6",
			want: []DiceGroup{
				{Count: 4, Sides: 6},
				{Count: 2, Sides: 10},
				{Count: 1, Sides: 6},
				{Count: 1, Sides: 6},
			},
		},
		{
			name:       "whitespace separated groups",
			expression: "4d6 d8\t2d12",
			want: []DiceGroup{
				{Count: 4, Sides: 6},
				{Count: 1, Sides: 8},
				{Count: 2, Sides: 12},
			},
		},
		{
			name:       "all supported die sizes",
			expression: "d2,d4,d6,d8,d10,d12,d20,d100",
			want: []DiceGroup{
				{Count: 1, Sides: 2},
				{Count: 1, Sides: 4},
				{Count: 1, Sides: 6},
				{Count: 1, Sides: 8},
				{Count: 1, Sides: 10},
				{Count: 1, Sides: 12},
				{Count: 1, Sides: 20},
				{Count: 1, Sides: 100},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseExpression(test.expression)
			if err != nil {
				t.Fatalf("ParseExpression() error = %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("ParseExpression() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestParseExpressionInvalid(t *testing.T) {
	tests := []string{
		"",
		"   ",
		",d6",
		"d6,",
		"d6,,d8",
		"d6, ,d8",
		"4 d6",
		"0d6",
		"d0",
		"2d3",
		"4d6+2",
		"d6 advantage",
		"d6/2",
		"1.5d6",
		"-1d6",
	}

	for _, expression := range tests {
		t.Run(expression, func(t *testing.T) {
			if groups, err := ParseExpression(expression); err == nil {
				t.Fatalf("ParseExpression(%q) = %#v, want an error", expression, groups)
			}
		})
	}
}

func TestParseExpressionLimits(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantErr    error
	}{
		{
			name:       "too many groups",
			expression: strings.TrimSuffix(strings.Repeat("d6,", MaxExpressionGroups+1), ","),
			wantErr:    ErrTooManyGroups,
		},
		{
			name:       "too many dice in one group",
			expression: "101d6",
			wantErr:    ErrTooManyDice,
		},
		{
			name:       "too many dice total",
			expression: strings.TrimSuffix(strings.Repeat("100d6,", MaxDicePerRequest/MaxDicePerGroup+1), ","),
			wantErr:    ErrTooManyDice,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseExpression(test.expression)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("ParseExpression() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestDiceGroupString(t *testing.T) {
	group := DiceGroup{Count: 4, Sides: 6}
	if got, want := group.String(), "4d6"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
