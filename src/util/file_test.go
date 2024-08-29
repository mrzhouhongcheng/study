package util

import (
	"fmt"
	"os"
	"testing"
)

func TestIsFile(t *testing.T) {
	path := "../../test/test.js"
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

func TestFileExists(t *testing.T) {
	path := "../../test/test.js"
	exist := FileExists(path)
	if !exist {
		t.Fatal("预计返回true, 结果返回false")
	}
	path = "../../test/test"
	exist = FileExists(path)
	if exist {
		t.Fatal("预计返回false, 结果返回true")
	}
}

// 文件复制测试
func TestCopyFile(t *testing.T) {
	path := "../../test/test.js"
	targetPath := "../../test/test_copy.js"
	err := CopyFile(path, targetPath)
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(targetPath)
}
