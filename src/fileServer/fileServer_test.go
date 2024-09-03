package fileserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	res, err := Split("../../test/test.js")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, len(res) > 0)
}

func TestSplit2(t *testing.T) {
	t.Skip("skipping Test TestSplit2 for now. because this jdb-8u361-windows-x64.exe file is not found")

	res, err := Split("../../test/jdk-8u361-windows-x64.exe")
	if err != nil {
		t.Error(err)
	}
	assert.True(t, len(res) > 0)
}

func TestSplitFilder(t *testing.T) {
	uuid, err := SplitFilder("../../test/test.js")

	if err != nil {
		t.Error(err)
	}
	assert.True(t, uuid != "")
}

func TestMerge(t *testing.T) {
	err := Merge("../../test/mergetest/8d616f67-dc25-41a7-a102-752d66aaffb7", "../../test/mergetest/8d616f67-dc25-41a7-a102-752d66aaffb7")
	assert.True(t, err == nil)
}
