package randomness

import (
	"errors"
	"math/rand"
)

var ErrSourceUnavailable = errors.New("randomness source is unavailable")

type Source struct {
	Description string
	Roll        func(sides int) (int, error)
}

type Table struct {
	PseudoRandom Source
}

func NewTable() Table {
	return Table{
		PseudoRandom: Source{
			Description: "Go math/rand pseudo-random generator (non-cryptographic)",
			Roll: func(sides int) (int, error) {
				return rand.Intn(sides) + 1, nil
			},
		},
	}
}
