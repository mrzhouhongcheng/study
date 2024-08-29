package util

import "testing"

func TestGenerageUUID(t *testing.T) {
	uuid := GenerageUUID()
	if uuid == "" {
		t.Error("generate UUID failed")
	}
}
