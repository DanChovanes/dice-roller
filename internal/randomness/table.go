package randomness

import (
	"errors"
	"math/rand"
)

var ErrSourceUnavailable = errors.New("randomness source is unavailable")

type Source func(sides int) (int, error)

type Table struct {
	PseudoRandom Source
}

func NewTable() Table {
	return Table{
		PseudoRandom: func(sides int) (int, error) {
			return rand.Intn(sides) + 1, nil
		},
	}
}
