package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	fileserver "com.zhouhc.study/src/fileServer"
	"com.zhouhc.study/src/util"
)

func main() {

	downJsonPath := flag.String("d", "", "指定数据down.json文件地址")
	checkJsonPath := flag.String("c", "", "指定检查的down.json文件地址")
	url := flag.String("h", "", "指定文件下载的url地址")
	output := flag.String("o", "", "指定文件的下载路径")
	flag.Parse()

	if *downJsonPath == "" && *checkJsonPath == "" && *url == "" {
		fmt.Println("Usage: down [option(-d|-c|-h)] <data.json>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *url != "" {
		err := downJSON(*url, *output)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		// 下载他的数据文件
		err = downPart(*output, *output)
		if err != nil {
			fmt.Printf("downPart failed: %v\n", err)
			os.Exit(1)
		}
		err = checkParts(*output)
		if err != nil {
			fmt.Printf("checkParts failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *downJsonPath != "" {
		err := downPart(*output, *output)
		if err != nil {
			fmt.Printf("downPart failed: %v\n", err)
			os.Exit(1)
		}
		err = checkParts(*output)
		if err != nil {
			fmt.Printf("checkParts failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *checkJsonPath != "" {
		err := checkParts(*output)
		if err != nil {
			fmt.Printf("checkParts failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

}

func downJSON(url, output string) error {
	params := make(map[string]string)
	params["filePath"] = url

	params_str, err := json.Marshal(params)
	if err != nil {
		fmt.Errorf("Unable to marshal params: %v", err)
		return err
	}
	res, err := http.Post("http://localhost:8889/down", "application/json", strings.NewReader(string(params_str)))
	if err != nil {
		fmt.Printf("Post request failed: %s\n", err)
		return err
	}
	defer res.Body.Close()
	dw, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Read response failed: %s\n", err)
		return err
	}
	return util.WriteFile(filepath.Join(output, "down.json"), dw)
}

func downPart(dwPath, output string) error {
	downJson, err := fileserver.GetDownjsonByPath(dwPath)
	if err != nil {
		return err
	}
	for _, path := range downJson.FileList {
		url := fmt.Sprintf("http://localhost:8888/%s", path)
		getFileByUrl(url, filepath.Join(output, filepath.Base(path)))
	}
	return nil
}

// dwPath : down.json 文件夹的路径地址
func checkParts(dwPath string) error {
	downJson, err := fileserver.GetDownjsonByPath(dwPath)
	if err != nil {
		fmt.Printf("Get down.json failed: %v\n", err)
		return err
	}
	for _, path := range downJson.FileList {
		// 获取文件在本地的路径
		part_path := filepath.Join(dwPath, filepath.Base(path))
		if !util.FileExists(part_path) {
			url := fmt.Sprintf("http://localhost:8888/%s", path)
			err = getFileByUrl(url, part_path)
			if err != nil {
				fmt.Println("Error getting file from URL: ", err)
				return err
			}
		}
	}
	// 合并文件
	err = fileserver.Merge(dwPath)
	if err != nil {
		fmt.Printf("Merge file failed: %v\n", err)
		return err
	}
	// 校验文件的hashkey
	code, err := util.CalculateFileHash(filepath.Join(dwPath, downJson.FileName))
	if err != nil {
		fmt.Println("Error calculating file hash from file: ", err)
		return err
	}
	if code != downJson.HashKey {
		return errors.New("hash key not supported")
	}
	return nil
}

func getFileByUrl(url, output string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return util.WriteFile(output, content)
}
