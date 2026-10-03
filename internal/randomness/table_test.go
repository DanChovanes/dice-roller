package randomness

import "testing"

func TestDefaultTableReturnsValidDieResult(t *testing.T) {
	table := NewTable()
	result, err := table.PseudoRandom(20)
	if err != nil {
		t.Fatalf("PseudoRandom() error = %v", err)
	}
	if result < 1 || result > 20 {
		t.Fatalf("PseudoRandom() = %d, want a result from 1 through 20", result)
	}
}
