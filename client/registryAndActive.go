package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	// 激活
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			go Active(host, port)
		}
	}()
}
