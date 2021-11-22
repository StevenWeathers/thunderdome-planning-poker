package db

import (
	"testing"
)

// TestHashApiKey calls hashApiKey and makes sure the return is not the same as the input
// and also calls hashApiKey twice with same key and makes sure the results match
func TestHashApiKey(t *testing.T) {
	KeyToHash := "infinitystones"
	HashedResult1 := hashApiKey(KeyToHash)
	HashedResult2 := hashApiKey(KeyToHash)
	HashedKey := "1af36dd2293e9a589fa31a72aaa894617eb662957914d98ee35d159fcf3fb7ab"

	if HashedResult1 == KeyToHash {
		t.Fatalf(`expected HashedResult1: %s to not match KeyToHash: %s`, HashedResult1, KeyToHash)
	}

	if HashedResult1 != HashedResult2 {
		t.Fatalf(`expected HashedResult1: %s to match HashedResult2: %s`, HashedResult1, HashedResult2)
	}

	if HashedResult1 != HashedKey {
		t.Fatalf(`expected HashedResult1: %s to match HashedKey: %s`, HashedResult1, HashedKey)
	}
}
