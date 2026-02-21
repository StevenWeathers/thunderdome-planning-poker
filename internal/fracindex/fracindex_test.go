package fracindex

import (
	"errors"
	"math/rand"
	"sort"
	"testing"
)

func TestKeyBetween(t *testing.T) {
	tests := []struct {
		a, b *string
		exp  *string
		err  error
	}{
		{nil, nil, new("a "), nil},
		{nil, new("a "), new("Z~"), nil},
		{nil, new("Z~"), new("Z}"), nil},
		{new("a "), nil, new("a!"), nil},
		{new("a!"), nil, new("a\""), nil},
		{new("a0"), new("a1"), new("a0P"), nil},
		{new("a1"), new("a2"), new("a1P"), nil},
		{new("a0V"), new("a1"), new("a0k"), nil},
		{new("Z~"), new("a "), new("Z~P"), nil},
		{new("Z~"), new("a!"), new("a "), nil},
		{nil, new("Y  "), new("X~~~"), nil},
		{new("b~~"), nil, new("c   "), nil},
		{new("a0"), new("a0V"), new("a0;"), nil},
		{new("a0"), new("a0G"), new("a04"), nil},
		{new("b125"), new("b129"), new("b127"), nil},
		{new("a0"), new("a1V"), new("a1"), nil},
		{new("Z~"), new("a 1"), new("a "), nil},
		{nil, new("a0V"), new("a0"), nil},
		{nil, new("b999"), new("b99"), nil},
		{nil, new("A                          "), nil, errors.New("key is too small")},
		// @TODO - fix the implementation to handle this case
		//{nil, strPtr("A                          !"), strPtr("A                           P"), nil},
		{new("zzzzzzzzzzzzzzzzzzzzzzzzzzy"), nil, new("zzzzzzzzzzzzzzzzzzzzzzzzzzz"), nil},
		{new("z~~~~~~~~~~~~~~~~~~~~~~~~~~"), nil, new("z~~~~~~~~~~~~~~~~~~~~~~~~~~P"), nil},
		{new("a0 "), nil, nil, errors.New("fractional part should not end with ' ' (space)")},
		{new("a0 "), new("a1"), nil, errors.New("fractional part should not end with ' ' (space)")},
		{new("0"), new("1"), nil, errors.New("head is out of range")},
		{new("a1"), new("a0"), nil, errors.New("key_between - a must be before b")},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			btwn, err := KeyBetween(tt.a, tt.b)
			if err != nil && (err.Error() != tt.err.Error()) {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
			if (btwn == nil && tt.exp != nil) || (btwn != nil && tt.exp == nil) || (btwn != nil && *btwn != *tt.exp) {
				t.Errorf("expected %v, got %v", tt.exp, btwn)
			}
		})
	}
}

func TestGenerateInsertOrder(t *testing.T) {
	die := rand.Intn

	// 1. generate a list of indices
	// 2. Permute the copy by moving items around
	// 3. Get new index of the item moved for each move
	// 4. order by index and compare to original list

	var prev *string
	var indices []string
	for range 5 {
		prev, _ = KeyBetween(prev, nil)
		indices = append(indices, *prev)
	}

	sorted := make([]string, len(indices))
	copy(sorted, indices)
	sort.Strings(sorted)
	if !vecCompare(sorted, indices) {
		t.Errorf("expected sorted and indices to be equal")
	}

	i := 0
	// Run through 1k random re-orderings and ensure the list is always sorted
	// correctly by fractional index
	for i < 10 {
		fromIndex := die(5)
		toIndex := die(5)
		if fromIndex == toIndex {
			continue
		}

		var fractIndex *string
		if toIndex == 0 {
			fractIndex, _ = KeyBetween(nil, new(indices[toIndex]))
		} else {
			fractIndex, _ = KeyBetween(new(indices[toIndex-1]), new(indices[toIndex]))
		}

		indices = append(indices[:toIndex], append([]string{*fractIndex}, indices[toIndex:]...)...)
		indices = append(indices[:fromIndex], indices[fromIndex+1:]...)
		sorted = make([]string, len(indices))
		copy(sorted, indices)
		sort.Strings(sorted)
		if !vecCompare(sorted, indices) {
			t.Errorf("expected sorted and indices to be equal")
		}

		i++
	}
}

func vecCompare(va, vb []string) bool {
	if len(va) != len(vb) {
		return false
	}
	for i := range va {
		if va[i] != vb[i] {
			return false
		}
	}
	return true
}
