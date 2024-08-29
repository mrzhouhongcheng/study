package fileserver

import "testing"


func TestMerage(t *testing.T) {
	err := Merge("../../test/test.js")
	if err != nil {
        t.Error(err)
    }
}


func TestMerage2(t *testing.T) {
	err := Merge("../../test/jdk-8u361-windows-x64.exe")
	if err != nil {
        t.Error(err)
    }
}

