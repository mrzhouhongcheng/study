package util

import (
	"fmt"
	"testing"
)

func TestIsFile(t *testing.T) {
	path := "../../test.js"
	res, err := IsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Error("预计返回true, 结果返回false")
	}
}

func TestIsDir(t *testing.T) {
	path := "../"
	res, err := IsDir(path)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Error("预计返回true, 结果返回false")
	}
}

func TestCalculateFileHash(t *testing.T) {
	path := "../../test/test.js"
	hash, err := CalculateFileHash(path)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)
}
