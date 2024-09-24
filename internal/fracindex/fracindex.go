// Package fracindex provides utilities for generating and manipulating
// lexicographically ordered keys. It allows for creating keys between existing keys,
// which is useful for maintaining sorted lists or implementing insertion operations
// in ordered data structures.
package fracindex

import (
	"errors"
	"strings"
)

const (
	Base95Digits = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	// SmallestInteger represents the smallest possible integer part of a key.
	SmallestInteger = "A                          "
	// IntegerZero represents the default starting key.
	IntegerZero = "a "
	// aCharcode is the ASCII code for the lowercase letter 'a'.
	aCharcode = 97
	// zCharcode is the ASCII code for the lowercase letter 'z'.
	zCharcode = 122
	// ACharcode is the ASCII code for the uppercase letter 'A'.
	ACharcode = 65
	// ZCharcode is the ASCII code for the uppercase letter 'Z'.
	ZCharcode = 90
	// minCharcode is the ASCII code for the space character.
	minCharcode = 32
)

// KeyBetween generates a key between two given keys a and b.
// If a is empty, it generates a key before b.
// If b is empty, it generates a key after a.
// If both are non-empty, it generates a key between a and b.
// Returns an error if the input keys are invalid or if a new key cannot be generated.
func KeyBetween(a, b *string) (*string, error) {
	digits := []rune(Base95Digits)

	if a != nil {
		if err := validateOrderKey(*a); err != nil {
			return nil, err
		}
	}
	if b != nil {
		if err := validateOrderKey(*b); err != nil {
			return nil, err
		}
	}

	switch {
	case a == nil && b == nil:
		zero := IntegerZero
		return &zero, nil
	case a != nil && b != nil:
		if *a > *b {
			return nil, errors.New("key_between - a must be before b")
		}

		ia, err := getIntegerPart(*a)
		if err != nil {
			return nil, err
		}
		ib, err := getIntegerPart(*b)
		if err != nil {
			return nil, err
		}
		fa := (*a)[len(ia):]
		fb := (*b)[len(ib):]
		if ia == ib {
			mid, err := midpoint([]rune(fa), []rune(fb), digits)
			if err != nil {
				return nil, err
			}
			result := ia + string(mid)
			return &result, nil
		}

		i, err := incrementInteger([]rune(ia), digits)
		if err != nil {
			return nil, err
		}
		if i != nil && *i < *b {
			return i, nil
		}

		mid, err := midpoint([]rune(fa), nil, digits)
		if err != nil {
			return nil, err
		}
		result := ia + string(mid)
		return &result, nil
	case a == nil:
		ib, err := getIntegerPart(*b)
		if err != nil {
			return nil, err
		}
		fb := (*b)[len(ib):]
		if ib == SmallestInteger {
			mid, err := midpoint([]rune(""), []rune(fb), digits)
			if err != nil {
				return nil, err
			}
			result := ib + string(mid)
			return &result, nil
		}
		if ib < *b {
			return &ib, nil
		}
		res, err := decrementInteger([]rune(ib), digits)
		if err != nil {
			return nil, err
		}
		return res, nil
	case b == nil:
		ia, err := getIntegerPart(*a)
		if err != nil {
			return nil, err
		}
		fa := (*a)[len(ia):]
		i, err := incrementInteger([]rune(ia), digits)
		if err != nil {
			return nil, err
		}
		if i == nil {
			mid, err := midpoint([]rune(fa), nil, digits)
			if err != nil {
				return nil, err
			}
			result := ia + string(mid)
			return &result, nil
		}
		return i, nil
	}
	return nil, nil
}

// midpoint calculates the midpoint between two strings a and b.
// If b is empty, it calculates the midpoint between a and the smallest possible value.
// The midpoint is determined lexicographically.
func midpoint(a, b []rune, digits []rune) ([]rune, error) {
	if b != nil {
		if string(a) == string(b) {
			return nil, errors.New("midpoint - a and b must not be equal")
		}
		if string(a) > string(b) {
			return nil, errors.New("midpoint - a must be before b")
		}
	}

	if len(a) > 0 && a[len(a)-1] == minCharcode || (b != nil && b[len(b)-1] == minCharcode) {
		return nil, errors.New("midpoint - a or b must not end with ' ' (space)")
	}

	if b != nil {
		n := 0
		for n < len(a) && a[n] == b[n] {
			n++
		}

		if n > 0 {
			mid, err := midpoint(a[n:], b, digits)
			if err != nil {
				return nil, err
			}
			return append(b[:n], mid...), nil
		}
	}

	var digitA int
	if len(a) > 0 {
		digitA = strings.IndexRune(string(digits), a[0])
	} else {
		digitA = 0
	}

	var digitB int
	if b != nil {
		digitB = strings.IndexRune(string(digits), b[0])
	} else {
		digitB = len(digits)
	}

	if digitB-digitA > 1 {
		midDigit := round(0.5 * float64(digitA+digitB))
		return []rune{digits[midDigit]}, nil
	} else {
		if len(b) > 1 {
			return b[:1], nil
		} else {
			mid, err := midpoint(a[1:], nil, digits)
			if err != nil {
				return nil, err
			}
			return append([]rune{digits[digitA]}, mid...), nil
		}
	}
}

// round rounds a float64 to the nearest integer.
// Returns the rounded integer.
func round(d float64) int {
	tenx := int(d * 10.0)
	truncated := int(d)
	if tenx-truncated*10 >= 5 {
		return truncated + 1
	}
	return truncated
}

// validateOrderKey validates an order key.
// Returns an error if the key is invalid.
func validateOrderKey(key string) error {
	if key == SmallestInteger {
		return errors.New("Key is too small")
	}
	i, err := getIntegerPart(key)
	if err != nil {
		return err
	}
	f := key[len(i):]
	if len(f) > 0 && f[len(f)-1] == minCharcode {
		return errors.New("Fractional part should not end with ' ' (space)")
	}
	return nil
}

// getIntegerPart extracts the integer part of an order key.
// Returns the integer part and an error if the key is invalid.
func getIntegerPart(key string) (string, error) {
	integerPartLen, err := getIntegerLen(rune(key[0]))
	if err != nil {
		return "", err
	}
	if integerPartLen > len(key) {
		return "", errors.New("integer part of key is too short")
	}
	return key[:integerPartLen], nil
}

// getIntegerLen returns the length of the integer part of an order key based on its head character.
// Returns the length and an error if the head character is invalid.
func getIntegerLen(head rune) (int, error) {
	if head >= aCharcode && head <= zCharcode {
		return int(head - aCharcode + 2), nil
	} else if head >= ACharcode && head <= ZCharcode {
		return int(ZCharcode - head + 2), nil
	} else {
		return 0, errors.New("head is out of range")
	}
}

// validateInteger validates the integer part of an order key.
// Returns an error if the integer part is invalid.
func validateInteger(i string) error {
	integerLen, err := getIntegerLen(rune(i[0]))
	if err != nil {
		return err
	}
	if len(i) != integerLen {
		return errors.New("invalid integer part of order key")
	}
	return nil
}

// incrementInteger increments the integer part of an order key.
// Returns the incremented integer part and an error if the operation fails.
func incrementInteger(x []rune, digits []rune) (*string, error) {
	if err := validateInteger(string(x)); err != nil {
		return nil, err
	}

	head := x[:1]
	digs := x[1:]
	carry := true

	for i := len(digs) - 1; i >= 0 && carry; i-- {
		temp := strings.IndexRune(string(digits), digs[i])
		if temp == -1 {
			return nil, errors.New("invalid digit")
		}
		d := temp + 1

		if d == len(digits) {
			digs[i] = digits[0]
		} else {
			digs[i] = digits[d]
			carry = false
		}
	}

	if carry {
		if string(head) == "Z" {
			zero := IntegerZero
			return &zero, nil
		}
		if string(head) == "z" {
			return nil, nil
		}
		h := head[0] + 1
		if h > aCharcode {
			digs = append(digs, digits[0])
		} else {
			digs = digs[:len(digs)-1]
		}
		result := string(append([]rune{h}, digs...))
		return &result, nil
	} else {
		result := string(append(head, digs...))
		return &result, nil
	}
}

// decrementInteger decrements the integer part of an order key.
// Returns the decremented integer part and an error if the operation fails.
func decrementInteger(x []rune, digits []rune) (*string, error) {
	if err := validateInteger(string(x)); err != nil {
		return nil, err
	}

	head := x[:1]
	digs := x[1:]
	borrow := true

	for i := len(digs) - 1; i >= 0 && borrow; i-- {
		temp := strings.IndexRune(string(digits), digs[i])
		if temp == -1 {
			return nil, errors.New("invalid digit")
		}
		d := temp - 1

		if d == -1 {
			digs[i] = digits[len(digits)-1]
		} else {
			digs[i] = digits[d]
			borrow = false
		}
	}

	if borrow {
		if string(head) == "a" {
			result := "Z" + string(digits[len(digits)-1])
			return &result, nil
		}
		if string(head) == "A" {
			return nil, nil
		}
		h := head[0] - 1
		if h < ZCharcode {
			digs = append(digs, digits[len(digits)-1])
		} else {
			digs = digs[:len(digs)-1]
		}
		result := string(append([]rune{h}, digs...))
		return &result, nil
	} else {
		result := string(append(head, digs...))
		return &result, nil
	}
}
