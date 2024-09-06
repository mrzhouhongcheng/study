package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"com.zhouhc.study/service"
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
	resp, err := client.Post("http://localhost:9999/gproxy/registry", "Content-type: application/json", bytes.NewBuffer(jsonData))
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
	resp, err := client.Post("http://localhost:9999/gproxy/active", "Content-type: application/json", bytes.NewBuffer(jsonData))
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

func InitializeGProxy() {
	err := Registry("localhost", 9912)
	if err != nil {
		log.Println("registry failed", err)
		panic(0)
	}
	// 激活
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			go Active("localhost", 9912)
		}
	}()
}
