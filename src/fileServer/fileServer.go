package fileserver

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"com.zhouhc.study/src/util"
)

func Merge(path string) ([]string, error) {
	// 如果路径不是一个文件, 那么就返回err
	isFile, err := util.IsFile(path)
	if err != nil {
		return nil, err
	}
	if !isFile {
		return nil, errors.New("path is not a file ")
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	file_name := fileInfo.Name()
	log.Println("file name is ", file_name)

	var buf = make([]byte, 1024*1024*10)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	index := 1
	var res []string = make([]string, 0)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("读取文件出错", err)
			return nil, err
		}
		if n <= 0 {
			break
		}
		newFilePath := file_name + ".part" + strconv.Itoa(index)
		newFilePath = filepath.Join(filepath.Dir(path), newFilePath)
		err = os.WriteFile(newFilePath, buf[:n], os.ModePerm)
		if err != nil {
			log.Println("write new file failed, ", err)
			return nil, err
		}
		index += 1
		res = append(res, newFilePath)
	}
	return res, nil
}

type DownJson struct {
	FileName   string   `json:"fileName"`
	FolderName string   `json:"folderName"`
	HashKey    string   `json:"hashKey"`
	FileList   []string `json:"fileList"`
}

func NewDownJons(uuid, hashKey, fileName string, fileList []string) *DownJson {
	return &DownJson{
		FileName:   fileName,
		FolderName: uuid,
		HashKey:    hashKey,
		FileList:   fileList,
	}
}

// 传入一个文件路径, 如果它是一个文件, 则创建一个对应的临时文件夹; 文件夹的名字是UUID生成的
// 复制这个文件到文件夹中, 对它进行分割;
// 将对应的文件夹路径uuid; 和对应的文件hashKey添加到一个JSON文件中
// 写入JSON文件, 文件名叫做down.json
// 需要下载文件 同时也需要写入到JSON文件当中;
func MergeFilder(path string) (string, error) {
	if !util.IsFileNotError(path) {
		log.Println("path is not a file")
		return "", errors.New("path is not a file, MergeFilder failed")
	}
	dirPath := filepath.Dir(path)
	uuid := util.GenerageUUID()
	targetPath := filepath.Join(dirPath, uuid, filepath.Base(path))
	util.CopyFile(path, targetPath)

	// 计算文件的hash值
	hashKey, _ := util.CalculateFileHash(targetPath)
	fileList, err := Merge(targetPath)
	if err != nil {
		return "", err
	}
	downjson := NewDownJons(uuid, hashKey, filepath.Base(targetPath), fileList)
	data, err := json.Marshal(downjson)
	if err != nil {
		return "", err
	}
	// 写入json到down.json中
	return uuid, util.WriteFile(filepath.Join(filepath.Dir(targetPath), "down.json"), data)
}

// 指定一个文件, 如果文件 > 10mb 那么就对这个文件进行切割;
// 如果文件不大于10mb; 那么就分割成一个文件
func Split(path string) error {
	return nil
}
