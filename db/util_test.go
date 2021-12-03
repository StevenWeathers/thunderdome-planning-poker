package db

import (
	"testing"
)

// TestHashString calls hashString and makes sure the return is not the same as the input
// and also calls hashString twice with same string and makes sure the results match
func TestHashString(t *testing.T) {
	StringToHash := "infinitystones"
	HashedResult1 := hashString(StringToHash)
	HashedResult2 := hashString(StringToHash)
	HashedString := "1af36dd2293e9a589fa31a72aaa894617eb662957914d98ee35d159fcf3fb7ab"

	if HashedResult1 == StringToHash {
		t.Fatalf(`expected HashedResult1: %s to not match StringToHash: %s`, HashedResult1, StringToHash)
	}

	if HashedResult1 != HashedResult2 {
		t.Fatalf(`expected HashedResult1: %s to match HashedResult2: %s`, HashedResult1, HashedResult2)
	}

	if HashedResult1 != HashedString {
		t.Fatalf(`expected HashedResult1: %s to match HashedString: %s`, HashedResult1, HashedString)
	}
}
