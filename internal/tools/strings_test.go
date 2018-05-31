package tools_test

import (
	"testing"

	"github.com/arsham/blush/internal/tools"
)

func TestIsPlainText(t *testing.T) {
	tcs := []struct {
		name  string
		input string
		want  bool
	}{
		{"null", string(0), true},
		{"space", " ", true},
		{"return", "\r", true},
		{"line feed", "\n", true},
		{"tab", "\t", true},
		{"bell", "\b", false},
		{"mix", "\n\n \r\nsjdk", true},
		{"1", "\x01", false},
		{"zero in middle", "n\x00b", true},
		{"bell in middle", "a\bc", false},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := tools.IsPlainText(tc.input)
			if got != tc.want {
				t.Errorf("tools.IsPlainText() = %t, want %t", tools.IsPlainText(tc.input), tc.want)
			}
		})
	}
}
