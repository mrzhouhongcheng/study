package util

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// 如果是文件则返回true
// 如果不是文件, 那么就返回false
func IsFile(path string) (bool, error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode().IsRegular(), nil
}

// 如果该路径是一个文件夹. 那么就返回true;
// 如果不是一个文件夹, 那么就返回false;
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// 计算文件的HASH码
// 如果成功, 那么就返回string, nil; 如果失败, 那么返回: "", error;
func CalculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)

	return fmt.Sprintf("%x", hash), nil
}
