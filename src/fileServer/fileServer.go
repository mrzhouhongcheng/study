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
	"sync"

	"com.zhouhc.study/src/util"
)

// path 源文件地址
// output 文件的输出目录
// Split is a function that splits a large file into smaller parts.//+
// It takes a single parameter://+
// - path: A string representing the path to the file to be split.//+
// //+
// The function returns two values://+
// - []string: A slice of strings representing the paths to the created part files.//+
// - error: An error value that is nil if the function completes successfully, or an error if an error occurs.//+
// //+
// The function reads the file specified by the given path and splits it into smaller parts.//+
// Each part file is named as the original file name followed by ".part" and a sequential index number.//+
// The part files are created in the same directory as the original file.//+
// The function uses a buffered reader to read the file and a wait group to ensure all part files are written before returning.//+
func Split(path, output string) ([]string, error) {
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
	syncChan := make(chan string, 20)
	var wg sync.WaitGroup
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
		newFilePath = filepath.Join(output, newFilePath)
		syncChan <- newFilePath
		wg.Add(1)
		go func(data []byte, filePath string) {
			log.Printf("split file : %s \n", filePath)
			defer func() {
				wg.Done()
				<-syncChan
			}()
			err := os.WriteFile(filePath, data, os.ModePerm)
			if err != nil {
				log.Println("write new file failed, ", err)
			}
		}(append([]byte(nil), buf[:n]...), newFilePath)
		index += 1
		res = append(res, newFilePath)
	}
	wg.Wait()
	return res, nil
}

type DownJson struct {
	FileName   string   `json:"fileName"`
	FolderName string   `json:"folderName"`
	FolderPath string   `json:"folderPath"`
	HashKey    string   `json:"hashKey"`
	FileList   []string `json:"fileList"`
}

func NewDownJons(uuid, hashKey, fileName, folderPath string, fileList []string) *DownJson {
	return &DownJson{
		FileName:   fileName,
		FolderName: uuid,
		HashKey:    hashKey,
		FolderPath: folderPath,
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
	uuid := util.GenerageUUID()
	targetPath := filepath.Join(os.TempDir(), uuid, filepath.Base(path))
	os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)

	// 计算文件的hash值
	log.Printf("calculating file: {%s} hash...\n", path)
	hashKey, _ := util.CalculateFileHash(path)
	log.Printf("current file hash key is %s \n", hashKey)
	fileList, err := Split(path, filepath.Dir(targetPath))
	if err != nil {
		return "", err
	}
	downjson := NewDownJons(uuid, hashKey, filepath.Base(targetPath), filepath.Dir(targetPath), fileList)
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
	if dwPath != "" && !util.IsFileNotError(dwPath) {
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
		log.Printf("merge file %v\n", val)
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
	data, err := os.ReadFile(path)
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
