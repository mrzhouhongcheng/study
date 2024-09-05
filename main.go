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
	"sync"
	"time"
)

type ProxyServer struct {
	Host       string    `json:"host"`
	Port       int       `json:"port"`
	LastAction time.Time `json:"lastAction"`
	IsAction   bool      `json:"isAction"`

	mu sync.Mutex
}

func NewProxyServerByModel(model ProxyServerModel) *ProxyServer {
	return &ProxyServer{
		Host: model.Host,
		Port: model.Port,
	}
}

type ProxyServerModel struct {
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	RequestTime time.Time `json:"requestTime"`
}

var ProxyServerMap map[string]*ProxyServer

func (p *ProxyServer) getKey() string {
	return p.Host + ":" + strconv.Itoa(p.Port)
}

// 注册
func registryHandler(w http.ResponseWriter, r *http.Request) {
	var model ProxyServerModel
	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	server := NewProxyServerByModel(model)
	key := server.getKey()
	if _, ok := ProxyServerMap[key]; ok {
		w.WriteHeader(http.StatusOK)
		return
	}
	server.mu.Lock()
	defer server.mu.Unlock()
	server.IsAction = false
	ProxyServerMap[key] = server

	w.WriteHeader(http.StatusOK)
}

// 激活
func activeHandler(w http.ResponseWriter, r *http.Request) {
	var model ProxyServerModel
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
	server := NewProxyServerByModel(model)
	key := server.getKey()
	if proxyServer, ok := ProxyServerMap[key]; ok {
		proxyServer.mu.Lock()
		defer proxyServer.mu.Unlock()
		proxyServer.LastAction = time.Now()
		w.WriteHeader(http.StatusOK)
	} else {
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
		proxy.mu.Lock()
		defer proxy.mu.Unlock()
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
	targetURL := "http://10.88.19.91"
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target url", http.StatusInternalServerError)
		return
	}
	var url string
	url = parsedURL.String() + r.URL.Path
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

	ProxyServerMap = make(map[string]*ProxyServer)
	http.HandleFunc("/", proxyHandler)
	http.HandleFunc("/active", activeHandler)
	http.HandleFunc("/registry", registryHandler)

	log.Println("Proxy server is running or port 9999....")
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
