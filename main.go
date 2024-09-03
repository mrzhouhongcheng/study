package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	fileserver "com.zhouhc.study/src/fileServer"
	"com.zhouhc.study/src/util"
)

// 如果需要下载, 第一步就是先对文件进行分割, 下载一个down.json文件, 这个文件的内容可以作为响应的结果进行返回
// 定义传入的参数:
// 传入的参数是这个http的下载地址. 但是两个是不同的路径;
// 需要启动两个服务, 一个是httpFIleServer, 一个是其他的关联程序, 两个监听的地址是不一样的
// 文件服务器监听的地址是8888
// 现在服务器监听的地址是8889
// 传入的参数就是文件服务器的下载地址

func startFTP() {
	go func() {
		fileServer := http.FileServer(http.Dir("./"))
		http.Handle("/", fileServer)
		if err := http.ListenAndServe(":8888", nil); err != nil {
			panic(err)
		}
	}()
}

func removeTempHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	var params = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	removePath, ok := params["removePath"]
	if !ok {
		http.Error(w, "Missing removePath parameter", http.StatusBadRequest)
		return
	}
	os.RemoveAll(removePath)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"code\":1,\"data\": \"success\"}"))
}

func downHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
	var params = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// 判断是否存在filePath这个请求参数
	filePath, ok := params["filePath"]
	if !ok {
		http.Error(w, "Missing filePath parameter", http.StatusBadRequest)
		return
	}
	// 对文件进行解析
	// 截取到:8888后面的字符串
	filePath = "." + filePath[strings.Index(filePath, ":8888")+5:]

	// 然后打印这个参数
	fmt.Println("filepath", filePath)
	if !util.IsFileNotError(filePath) {
		http.Error(w, "传入的地址不是一个文件路径", http.StatusBadRequest)
		return
	}
	uuid, err := fileserver.SplitFilder(filePath)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusInternalServerError)
		return
	}

	// 读取一个文件
	data, err := os.ReadFile(filepath.Join(os.TempDir(), uuid, "down.json"))
	if err != nil {
		http.Error(w, "read down json file is error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(data))
}

func downpartHandler(w http.ResponseWriter, r *http.Request) {
	part := r.URL.Query().Get("part")
	if !util.FileExists(part) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	FileName := filepath.Base(part)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", FileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, part)
}

func main() {
	startFTP()

	http.HandleFunc("/down", downHandler)
	http.HandleFunc("/remove", removeTempHandler)
	http.HandleFunc("/downpart", downpartHandler)

	log.Fatal(http.ListenAndServe(":8889", nil))

}
