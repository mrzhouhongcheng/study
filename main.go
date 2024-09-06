package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"com.zhouhc.study/service"
	"github.com/spf13/viper"
)

var ProxyServerMap map[string]*service.ProxyServer

// 注册
func registryHandler(w http.ResponseWriter, r *http.Request) {
	var model service.ProxyServerModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	server := service.NewProxyServerByModel(model)
	key := server.GetKey()
	if _, ok := ProxyServerMap[key]; ok {
		w.WriteHeader(http.StatusOK)
		return
	}
	server.Mu.Lock()
	defer server.Mu.Unlock()
	server.IsAction = false
	ProxyServerMap[key] = server

	w.WriteHeader(http.StatusOK)
}

func activeCheck() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if len(ProxyServerMap) == 0 {
				continue
			}
			for _, val := range ProxyServerMap {
				if !val.IsAction {
					continue
				}
				timeout := val.LastAction.Add(3 * 20 * time.Second)
				if time.Now().After(timeout) {
					val.Mu.Lock()
					val.IsAction = false
					val.Mu.Unlock()
				}
			}
		}
	}()
}

// 激活
func activeHandler(w http.ResponseWriter, r *http.Request) {
	var model service.ProxyServerModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	timeout := model.RequestTime.Add(20 * time.Second)
	if time.Now().After(timeout) {
		http.Error(w, "Invalid request time", http.StatusInternalServerError)
		return
	}
	server := service.NewProxyServerByModel(model)
	key := server.GetKey()
	if proxyServer, ok := ProxyServerMap[key]; ok {
		proxyServer.Mu.Lock()
		defer proxyServer.Mu.Unlock()
		proxyServer.LastAction = time.Now()
		proxyServer.IsAction = true
		w.WriteHeader(http.StatusOK)
	} else {
		server.Mu.Lock()
		defer server.Mu.Unlock()
		server.IsAction = false
		server.LastAction = time.Now()
		ProxyServerMap[key] = server
		http.Error(w, "proxy server not registry", http.StatusInternalServerError)
	}
}

// map中随机选择一个地址
func geturlByProxyServerMap() (string, error) {
	if len(ProxyServerMap) == 0 {
		return "", errors.New("proxy server map is empty")
	}
	keys := make([]string, 0)
	for key, proxy := range ProxyServerMap {
		proxy.Mu.Lock()
		defer proxy.Mu.Unlock()
		if proxy.IsAction {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return "", errors.New("no active proxy server")
	}
	randomIndex := rand.Intn(len(keys))
	proxyServer := ProxyServerMap[keys[randomIndex]]
	return fmt.Sprintf("http://%s:%d", proxyServer.Host, proxyServer.Port), nil
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {

	targetURL, err := geturlByProxyServerMap()
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target url", http.StatusInternalServerError)
		return
	}
	url := parsedURL.String() + r.URL.Path
	log.Println("proxy url : ", url)
	proxyReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}
	proxyReq.Header = r.Header

	client := &http.Client{}

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to connect to target serve", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	for key, Values := range resp.Header {
		for _, val := range Values {
			w.Header().Add(key, val)
		}
	}
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying response body", http.StatusInternalServerError)
		return
	}
}

func main() {
	viper.SetConfigName("application.yml")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	port := viper.GetInt("port")
	log.Println("gproxy port listen: ", port)

	ProxyServerMap = make(map[string]*service.ProxyServer)
	go activeCheck()

	http.HandleFunc("/", proxyHandler)
	http.HandleFunc("/gproxy/active", activeHandler)
	http.HandleFunc("/gproxy/registry", registryHandler)

	log.Println("Proxy server is running or port " + strconv.Itoa(port) + "....")
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
