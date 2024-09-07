package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"com.zhouhc.study/service"
	"github.com/spf13/viper"
)

func Registry(host string, port int) error {
	params := service.ProxyServerModel{
		Host:        host,
		Port:        port,
		RequestTime: time.Now(),
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		log.Println("to json string failed", err)
		return err
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(
		fmt.Sprintf("%sgproxy/registry", getServerURL()),
		"Content-type: application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("http post failed", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resData, _ := io.ReadAll(resp.Body)
		log.Println("http status code not 200: ", string(resData))
		return errors.New("HTTP status code not 200")
	}
	return nil
}

func Active(host string, port int) error {
	params := service.ProxyServerModel{
		Host: host,
		Port: port,

		RequestTime: time.Now(),
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		log.Println("to json string failed", err)
		return err
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(
		fmt.Sprintf("%sgproxy/active", getServerURL()),
		"Content-type: application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("http post failed", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resData, _ := io.ReadAll(resp.Body)
		log.Println("http status code not 200: ", string(resData))
		return errors.New("HTTP status code not 200")
	}
	return nil
}

func getServerURL() string {
	url := viper.GetString("gproxy.server")
	if url == "" {
		log.Fatal("gproxy.server is required")
		panic(0)
	}
	if strings.HasSuffix(url, "/") {
		return url
	}
	return url + "/"
}

func InitializeGProxy() {
	port := viper.GetInt("gproxy.port")
	host := viper.GetString("gproxy.host")
	err := Registry(host, port)
	if err != nil {
		log.Println("registry failed", err)
		panic(0)
	}
	Active(host, port)
	// 激活
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			Active(host, port)
		}
	}()

	// 程序退出时需要执行
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
        sig := <-sigChan
		fmt.Println("收到信号:", sig)
		fmt.Println("remove server")

		os.Exit(0)
	}()
}
