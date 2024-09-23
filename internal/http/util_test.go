package http

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestMain(m *testing.M) {
	validate = validator.New()
	m.Run()
}

func TestValidateUserAccountWithPasswords(t *testing.T) {
	tests := []struct {
		name           string
		inputName      string
		inputEmail     string
		inputPassword1 string
		inputPassword2 string
		wantErr        bool
	}{
		{
			name:           "Valid user account",
			inputName:      "Thor",
			inputEmail:     "thor@thunderdome.dev",
			inputPassword1: "lokiIsAJoke",
			inputPassword2: "lokiIsAJoke",
			wantErr:        false,
		},
		{
			name:           "Invalid email",
			inputName:      "Thor",
			inputEmail:     "thor",
			inputPassword1: "lokiIsAJoke",
			inputPassword2: "lokiIsAJoke",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, email, password, err := validateUserAccountWithPasswords(tt.inputName, tt.inputEmail, tt.inputPassword1, tt.inputPassword2)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateUserAccountWithPasswords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if name != tt.inputName || email != tt.inputEmail || password != tt.inputPassword1 {
					t.Errorf("validateUserAccountWithPasswords() = (%v, %v, %v), want (%v, %v, %v)",
						name, email, password, tt.inputName, tt.inputEmail, tt.inputPassword1)
				}
			}
		})
	}
}

func TestValidateUserPassword(t *testing.T) {
	tests := []struct {
		name      string
		password1 string
		password2 string
		wantErr   bool
	}{
		{
			name:      "Matching passwords",
			password1: "lokiIsAJoke",
			password2: "lokiIsAJoke",
			wantErr:   false,
		},
		{
			name:      "Non-matching passwords",
			password1: "lokiIsAJoke",
			password2: "lokiIsAJokeFail",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := validateUserPassword(tt.password1, tt.password2)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateUserPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && password != tt.password1 {
				t.Errorf("validateUserPassword() = %v, want %v", password, tt.password1)
			}
		})
	}
}

func TestGetLimitOffsetFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    url.Values
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "Default values",
			queryParams:    url.Values{},
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "Valid limit and offset",
			queryParams:    url.Values{"limit": []string{"30"}, "offset": []string{"10"}},
			expectedLimit:  30,
			expectedOffset: 10,
		},
		{
			name:           "Invalid limit (use default)",
			queryParams:    url.Values{"limit": []string{"invalid"}, "offset": []string{"5"}},
			expectedLimit:  20,
			expectedOffset: 5,
		},
		{
			name:           "Invalid offset (use default)",
			queryParams:    url.Values{"limit": []string{"25"}, "offset": []string{"invalid"}},
			expectedLimit:  25,
			expectedOffset: 0,
		},
		{
			name:           "Zero limit (use default)",
			queryParams:    url.Values{"limit": []string{"0"}, "offset": []string{"15"}},
			expectedLimit:  20,
			expectedOffset: 15,
		},
		{
			name:           "Negative values (use defaults)",
			queryParams:    url.Values{"limit": []string{"-10"}, "offset": []string{"-5"}},
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "Negative limit, valid offset",
			queryParams:    url.Values{"limit": []string{"-10"}, "offset": []string{"5"}},
			expectedLimit:  20,
			expectedOffset: 5,
		},
		{
			name:           "Valid limit, negative offset",
			queryParams:    url.Values{"limit": []string{"30"}, "offset": []string{"-5"}},
			expectedLimit:  30,
			expectedOffset: 0,
		},
		{
			name:           "Very large values",
			queryParams:    url.Values{"limit": []string{"1000000"}, "offset": []string{"9999999"}},
			expectedLimit:  1000000,
			expectedOffset: 9999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/?"+tt.queryParams.Encode(), nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			limit, offset := getLimitOffsetFromRequest(req)

			if limit != tt.expectedLimit {
				t.Errorf("Expected limit %d, but got %d", tt.expectedLimit, limit)
			}
			if offset != tt.expectedOffset {
				t.Errorf("Expected offset %d, but got %d", tt.expectedOffset, offset)
			}
		})
	}
}

func TestGetSearchFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    url.Values
		expectedSearch string
		expectError    bool
	}{
		{
			name:           "Valid search query",
			queryParams:    url.Values{"search": []string{"validquery"}},
			expectedSearch: "validquery",
			expectError:    false,
		},
		{
			name:           "Minimum length search query",
			queryParams:    url.Values{"search": []string{"abc"}},
			expectedSearch: "abc",
			expectError:    false,
		},
		{
			name:           "Empty search query",
			queryParams:    url.Values{"search": []string{""}},
			expectedSearch: "",
			expectError:    true,
		},
		{
			name:           "Missing search query",
			queryParams:    url.Values{},
			expectedSearch: "",
			expectError:    true,
		},
		{
			name:           "Too short search query",
			queryParams:    url.Values{"search": []string{"ab"}},
			expectedSearch: "",
			expectError:    true,
		},
		{
			name:           "Long search query",
			queryParams:    url.Values{"search": []string{"thisisaverylongsearchquery"}},
			expectedSearch: "thisisaverylongsearchquery",
			expectError:    false,
		},
		{
			name:           "Search query with spaces",
			queryParams:    url.Values{"search": []string{"search with spaces"}},
			expectedSearch: "search with spaces",
			expectError:    false,
		},
		{
			name:           "Search query with special characters",
			queryParams:    url.Values{"search": []string{"search@#$%^&*"}},
			expectedSearch: "search@#$%^&*",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/?"+tt.queryParams.Encode(), nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			search, err := getSearchFromRequest(req)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if search != tt.expectedSearch {
					t.Errorf("Expected search '%s', but got '%s'", tt.expectedSearch, search)
				}
			}
		})
	}
}

func TestSanitizeUserInputForLogs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No newlines",
			input:    "This is a normal string",
			expected: "This is a normal string",
		},
		{
			name:     "With Unix newlines",
			input:    "This has\na newline",
			expected: "This hasa newline",
		},
		{
			name:     "With Windows newlines",
			input:    "This has\r\na Windows newline",
			expected: "This hasa Windows newline",
		},
		{
			name:     "With mixed newlines",
			input:    "This has\nmixed\r\nnewlines",
			expected: "This hasmixednewlines",
		},
		{
			name:     "Only newlines",
			input:    "\n\r\n\n\r\n",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String with only spaces",
			input:    "   ",
			expected: "   ",
		},
		{
			name:     "Newlines at start and end",
			input:    "\n\r\nThis is a test\n\r\n",
			expected: "This is a test",
		},
		{
			name:     "With tabs",
			input:    "This\thas\ttabs\nand newlines",
			expected: "This\thas\ttabsand newlines",
		},
		{
			name:     "With other whitespace characters",
			input:    "This has\u000B\u000C",
			expected: "This has\u000B\u000C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeUserInputForLogs(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', but got '%s'", tt.expected, result)
			}

			// Additional check to ensure no newlines remain
			if strings.Contains(result, "\n") || strings.Contains(result, "\r") {
				t.Errorf("Sanitized string still contains newline characters: %q", result)
			}
		})
	}
}
