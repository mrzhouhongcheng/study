package fileserver

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"com.zhouhc.study/src/util"
)

func Merge(path string) error {
	// 如果路径不是一个文件, 那么就返回err
	isFile, err := util.IsFile(path)
	if err != nil {
		return err
	}
	if !isFile {
		return errors.New("path is not a file ")
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	file_name := fileInfo.Name()
	log.Println("file name is ", file_name)

	var buf = make([]byte, 1024*1024*10)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	index := 1
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("读取文件出错", err)
			return err
		}
		if n <= 0 {
			break
		}
		newFilePath := file_name + ".part" + strconv.Itoa(index)
		newFilePath = filepath.Join(filepath.Dir(path), newFilePath)
		err = os.WriteFile(newFilePath, buf[:n], os.ModePerm)
		if err != nil {
			log.Println("write new file failed, ", err)
			return err
		}
		index += 1
	}
	return nil
}

// 指定一个文件, 如果文件 > 10mb 那么就对这个文件进行切割;
// 如果文件不大于10mb; 那么就分割成一个文件
func Split(path string) error {
	return nil
}
