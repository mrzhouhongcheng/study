package util

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

func IsFileNotError(path string) bool {
	isFile, _ := IsFile(path)
	return isFile
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

func IsDirNotError(path string) bool {
	isDir, _ := IsDir(path)
	return isDir
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

// FileExists判断文件是否存在;
// 如果存在则返回true; 否则返回false
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// 文件复制, 传入两个文件路径
// 如果sourcePath不存在, 那么就报错
// 如果targetPath已经存在, 则先删除后新建
// 如果targetPath不存在, 那么就创建一个文件
func CopyFile(sourcePath, targetPath string) error {
	if !FileExists(sourcePath) || !IsFileNotError(sourcePath) {
		log.Println("复制文件出错; 文件路径" + sourcePath + "不存在或者不是一个文件")
		return errors.New("copy file failed")
	}
	if FileExists(targetPath) {
		err := os.Remove(targetPath)
		if err != nil {
			log.Println("删除旧文件出错; ", err)
			return err
		}
	}
	dirPath := filepath.Dir(targetPath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Println("创建文件夹出错; ", err)
		return err
	}
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Println("open source file failed; ", err)
		return err
	}
	defer sourceFile.Close()
	targetFile, err := os.Create(targetPath)
	if err != nil {
		log.Println("create target file failed; ", err)
		return err
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		log.Println("copy file failed; ", err)
		return err
	}
	return nil
}

// 写入一个文件, 如果文件存在, 那么就删除这个文件,
// 如果文件不存在, 那么就创建一个文件
func WriteFile(filePath string, data []byte) error {
	// 判断文件是否存在
	if FileExists(filePath) {
		os.Remove(filePath)
	}
	// 获取他的父目录地址
	dirPath := filepath.Dir(filePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}
