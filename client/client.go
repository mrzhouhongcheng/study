package main

import (
	"encoding/json"
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
			fmt.Errorf("%v", err)
			os.Exit(1)
		}
		// 下载他的数据文件
		err = downPart(*output, *output)
		if err != nil {
			fmt.Printf("downPart failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *downJsonPath != "" {
		fmt.Println("-d is ", *downJsonPath)
	}

	if *checkJsonPath != "" {
		fmt.Println("-c is ", *checkJsonPath)
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
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		util.WriteFile(filepath.Join(output, filepath.Base(path)), content)
	}
	return nil
}
