package thunderdome

import (
	"strings"
	"testing"
)

func TestNormalizeRetroTag(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "unchanged", input: "Blocked", want: "Blocked"},
		{name: "collapses whitespace", input: "  Not \t team\n  responsibility ", want: "Not team responsibility"},
		{name: "empty", input: "", wantErr: true},
		{name: "only whitespace", input: " \n\t ", wantErr: true},
		{name: "max length counts characters", input: strings.Repeat("é", RetroTagMaxLength), want: strings.Repeat("é", RetroTagMaxLength)},
		{name: "too long", input: strings.Repeat("a", RetroTagMaxLength+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeRetroTag(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizeRetroTag(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("NormalizeRetroTag(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
