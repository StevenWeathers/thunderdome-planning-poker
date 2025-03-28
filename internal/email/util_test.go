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
		// Test 9: Handling an invalid error case (although unlikely)
		{
			// Here we are assuming there might be some edge case or failure.
			// For now, we expect no errors unless transform implementation changes.
			input:    string([]byte{0x80}), // Invalid byte sequence
			expected: "",
			err:      true,
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
