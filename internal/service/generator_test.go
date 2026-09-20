package service

import "testing"

func TestGenerateShortURL(t *testing.T) {
	for i := 0; i < 1000; i++ {
		got, err := GenerateShortURL()
		if err != nil {
			t.Fatalf("GenerateShortURL() returned error: %v", err)
		}

		if len(got) != shortURLLength {
			t.Fatalf("length = %d, want %d", len(got), shortURLLength)
		}

		for _, char := range got {
			if !isAllowedCharacter(char) {
				t.Fatalf("generated invalid character: %q", char)
			}
		}
	}
}

func TestIsValidShortURL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "valid",
			value: "aB123_XyZ9",
			want:  true,
		},
		{
			name:  "too short",
			value: "abc123",
			want:  false,
		},
		{
			name:  "too long",
			value: "abc123_XYZ_extra",
			want:  false,
		},
		{
			name:  "invalid character",
			value: "abc123-XYZ",
			want:  false,
		},
		{
			name:  "space",
			value: "abc123 XYZ",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidShortURL(tt.value)

			if got != tt.want {
				t.Fatalf(
					"IsValidShortURL(%q) = %v, want %v",
					tt.value,
					got,
					tt.want,
				)
			}
		})
	}
}
