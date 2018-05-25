package tools_test

import (
	"github.com/arsham/blush/internal/tools"
)

func ExampleFiles() {
	tools.Files(true, "~/Documents", "/tmp")
	// Or
	dirs := []string{"~/Documents", "/tmp"}
	tools.Files(false, dirs...)
}
