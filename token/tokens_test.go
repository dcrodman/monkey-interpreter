package token

import "testing"

func TestIsValidIdentifier(t *testing.T) {
	tests := []struct {
		literal string
		want    bool
	}{
		{"x", true},
		{"is_valid", true},
		{"_valid", true},
		{"valid_", true},
		{"abc123", true},
		{"1abc", false},
		{"abc 123", false},
	}
	for _, tt := range tests {
		if got := IsValidIdentifier(tt.literal); got != tt.want {
			t.Errorf("IsValidIdentifier() = %v, want %v", got, tt.want)
		}
	}
}

func TestIsValidInteger(t *testing.T) {
	tests := []struct {
		literal string
		want    bool
	}{
		{"1", true},
		{"12345", true},
		{"9_01", false},
		{"a", false},
	}
	for _, tt := range tests {
		if got := IsValidInteger(tt.literal); got != tt.want {
			t.Errorf("IsValidInteger() = %v, want %v", got, tt.want)
		}
	}
}

func TestGetKeywordType(t *testing.T) {
	tests := []struct {
		literal string
		want    bool
		want1   TokenType
	}{
		{"fn", true, FUNCTION},
		{"let", true, LET},
		{"if", true, IF},
		{"else", true, ELSE},
		{"return", true, RETURN},
		{"true", true, TRUE},
		{"false", true, FALSE},
		{"fail", false, ""},
		{"", false, ""},
	}
	for _, tt := range tests {
		got, got1 := GetKeywordType(tt.literal)
		if got != tt.want {
			t.Errorf("GetKeywordType() got = %v, want %v", got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("GetKeywordType() got1 = %v, want %v", got1, tt.want1)
		}
	}
}
