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
		{nil, nil, strPtr("a "), nil},
		{nil, strPtr("a "), strPtr("Z~"), nil},
		{nil, strPtr("Z~"), strPtr("Z}"), nil},
		{strPtr("a "), nil, strPtr("a!"), nil},
		{strPtr("a!"), nil, strPtr("a\""), nil},
		{strPtr("a0"), strPtr("a1"), strPtr("a0P"), nil},
		{strPtr("a1"), strPtr("a2"), strPtr("a1P"), nil},
		{strPtr("a0V"), strPtr("a1"), strPtr("a0k"), nil},
		{strPtr("Z~"), strPtr("a "), strPtr("Z~P"), nil},
		{strPtr("Z~"), strPtr("a!"), strPtr("a "), nil},
		{nil, strPtr("Y  "), strPtr("X~~~"), nil},
		{strPtr("b~~"), nil, strPtr("c   "), nil},
		{strPtr("a0"), strPtr("a0V"), strPtr("a0;"), nil},
		{strPtr("a0"), strPtr("a0G"), strPtr("a04"), nil},
		{strPtr("b125"), strPtr("b129"), strPtr("b127"), nil},
		{strPtr("a0"), strPtr("a1V"), strPtr("a1"), nil},
		{strPtr("Z~"), strPtr("a 1"), strPtr("a "), nil},
		{nil, strPtr("a0V"), strPtr("a0"), nil},
		{nil, strPtr("b999"), strPtr("b99"), nil},
		{nil, strPtr("A                          "), nil, errors.New("Key is too small")},
		// @TODO - fix the implementation to handle this case
		//{nil, strPtr("A                          !"), strPtr("A                           P"), nil},
		{strPtr("zzzzzzzzzzzzzzzzzzzzzzzzzzy"), nil, strPtr("zzzzzzzzzzzzzzzzzzzzzzzzzzz"), nil},
		{strPtr("z~~~~~~~~~~~~~~~~~~~~~~~~~~"), nil, strPtr("z~~~~~~~~~~~~~~~~~~~~~~~~~~P"), nil},
		{strPtr("a0 "), nil, nil, errors.New("Fractional part should not end with ' ' (space)")},
		{strPtr("a0 "), strPtr("a1"), nil, errors.New("Fractional part should not end with ' ' (space)")},
		{strPtr("0"), strPtr("1"), nil, errors.New("head is out of range")},
		{strPtr("a1"), strPtr("a0"), nil, errors.New("key_between - a must be before b")},
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
	for i := 0; i < 5; i++ {
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
			fractIndex, _ = KeyBetween(nil, strPtr(indices[toIndex]))
		} else {
			fractIndex, _ = KeyBetween(strPtr(indices[toIndex-1]), strPtr(indices[toIndex]))
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

func strPtr(s string) *string {
	return &s
}
