package tools_test

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/internal/tools"
)

func TestIsPlainText(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name  string
		input string
		want  bool
	}{
		{"null", fmt.Sprintf("%d", 0), true},
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tools.IsPlainText(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
