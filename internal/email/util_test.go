package email

import (
	"testing"
)

// Test for removeAccents function
func TestRemoveAccents(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		// Test 1: String without accents
		{
			input:    "hello",
			expected: "hello",
			err:      false,
		},
		// Test 2: String with accents (common case)
		{
			input:    "héllo",
			expected: "hello",
			err:      false,
		},
		// Test 3: String with multiple accents
		{
			input:    "façade",
			expected: "facade",
			err:      false,
		},
		// Test 4: String with combined accented characters
		{
			input:    "café",
			expected: "cafe",
			err:      false,
		},
		// Test 5: Empty string
		{
			input:    "",
			expected: "",
			err:      false,
		},
		// Test 6: String with only non-accented characters
		{
			input:    "hello world!",
			expected: "hello world!",
			err:      false,
		},
		// Test 7: String with special characters (no accents)
		{
			input:    "123!@#$",
			expected: "123!@#$",
			err:      false,
		},
		// Test 8: String with a mix of accented and non-accented characters
		{
			input:    "niño, jalapeño",
			expected: "nino, jalapeno",
			err:      false,
		},
		// Test 9: Handling an invalid byte sequence
		{
			// Input with an invalid byte sequence (0x80 is not valid UTF-8)
			input:    string([]byte{0x80}),
			expected: "�",   // The replacement character that Go uses for invalid UTF-8 sequences
			err:      false, // No error is expected, just replacement
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output, err := removeAccents(test.input)

			// Check if we expect an error
			if test.err && err == nil {
				t.Errorf("expected error, but got nil")
			}

			// Check if we don't expect an error
			if !test.err && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check the output
			if output != test.expected {
				t.Errorf("expected %v, got %v", test.expected, output)
			}
		})
	}
}
