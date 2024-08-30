package fileserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	res, err := Merge("../../test/test.js")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, len(res) > 0)
}

func TestMerge2(t *testing.T) {
	t.Skip("skipping Test TestMerge2 for now. because this jdb-8u361-windows-x64.exe file is not found")

	res, err := Merge("../../test/jdk-8u361-windows-x64.exe")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, len(res) > 0)
}

func TestMergeFilder(t *testing.T) {
	uuid, err := MergeFilder("../../test/test.js")

	if err != nil {
		t.Error(err)
	}
	assert.True(t, uuid != "")
}
