package util

import "os"

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
