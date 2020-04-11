package dechoconf

import (
	"testing"
)

func TestDecodeFile(t *testing.T) {
	err := DecodeFile("application.properties")
	if err.Error() != "application.properties doesn't have toml / yaml file extensions" {
		t.Error("Failed to detect unsupported file: application.properties")
	}
}
