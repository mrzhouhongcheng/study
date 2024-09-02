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

	"com.zhouhc.study/src/util"
)

func main() {

	downJsonPath := flag.String("d", "", "指定数据down.json文件地址")
	checkJsonPath := flag.String("c", "", "指定检查的down.json文件地址")
	url := flag.String("h", "", "指定文件下载的url地址")
	path := flag.String("o", "", "指定文件的下载路径")
	fmt.Println("path", path)
	flag.Parse()

	if *downJsonPath == "" && *checkJsonPath == "" && *url == "" {
		fmt.Println("Usage: down [option(-d|-c|-h)] <data.json>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *url != "" {
		params := make(map[string]string)
		params["filePath"] = *url

		params_str, err := json.Marshal(params)
		if err != nil {
			fmt.Errorf("Unable to marshal params: %v", err)
			os.Exit(1)
		}
		res, err := http.Post(*url, "application/json", strings.NewReader(string(params_str)))
		if err != nil {
			fmt.Printf("Post request failed: %s\n", err)
			os.Exit(1)
		}
		defer res.Body.Close()
		dw, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Read response failed: %s\n", err)
			os.Exit(1)
		}
		// 将内容写入到指定的位置

		fmt.Println("-h is ", *url)
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
	if output != "" {
		util.WriteFile(output, dw)
	} else {
		util.WriteFile("./down.json", dw)
	}
	return nil
}
