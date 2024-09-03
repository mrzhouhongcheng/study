package fileserver

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"com.zhouhc.study/src/util"
)

func Split(path string) ([]string, error) {
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
func SplitFilder(path string) (string, error) {
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
	fileList, err := Split(targetPath)
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

// 合并文件
// 传入一个文件夹路径, 然后读取这个文件夹中的down.json文件
// 根据down.json文件, 创建一个fileName文件
// 根据down.json文件中的fileList, 依次将内容写入到fileName文件中
// 校验文件的hashKey编码
// 如果传入的不是一个文件夹, 则报错
// 如果文件夹路径下没有down.json文件, 则报错
// 如果没有找到fileList中相对应的分片文件, 则报错
// 如果hashKey不匹配, 则报错
func Merge(dwPath, output string) error {
	// 判断传入的路径是否是一个文件夹
	if dwPath != "" && !util.IsDirNotError(dwPath) {
		return errors.New("path is not a directory")
	}
	// 读取down.json
	downJson, err := GetDownjsonByPath(dwPath)
	if err != nil {
		return err
	}
	filePath := filepath.Join(output, downJson.FileName)
	if util.FileExists(filePath) {
		os.Remove(filePath)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, val := range downJson.FileList {
		partName := filepath.Base(val)
		partFile, err := os.Open(filepath.Join(output, partName))
		if err != nil {
			return err
		}
		_, err = io.Copy(file, partFile)
		if err != nil {
			return err
		}
		partFile.Close()
	}
	return nil
}

func GetDownjsonByPath(path string) (*DownJson, error) {

	data, err := os.ReadFile(filepath.Join(path, "down.json"))
	if err != nil {
		return nil, err
	}
	var res DownJson
	return &res, json.Unmarshal(data, &res)
}

// 文件下载.
func DownPartFile(partFileUrl, outPath string) error {
	res, err := http.Get(partFileUrl)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	return util.WriteFile(outPath, data)
}
