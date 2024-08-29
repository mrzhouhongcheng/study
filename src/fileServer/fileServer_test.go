package fileserver

import "testing"

func TestMerge(t *testing.T) {
	err := Merge("../../test/test.js")
	if err != nil {
		t.Error(err)
	}
}

func TestMerge2(t *testing.T) {
	t.Skip("skipping Test TestMerge2 for now. because this jdb-8u361-windows-x64.exe file is not found")

	err := Merge("../../test/jdk-8u361-windows-x64.exe")
	if err != nil {
		t.Error(err)
	}
}

func TestMergeFilder(t *testing.T) {
	t.Skip("skipping Test TestMergeFilder for now. because this folder does not exist")
	err := MergeFilder("../../test/jdk-8u361-windows-x64.exe")

	if err != nil {
		t.Error(err)
	}
}
