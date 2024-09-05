package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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

var ProxyServerMap map[string]ProxyServer 

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := "http://10.88.19.91"
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target url", http.StatusInternalServerError)
		return
	}
	var url string
	if strings.HasPrefix(r.URL.Path, "/api") {
		url = parsedURL.String() + r.URL.Path
	} else {
		url = parsedURL.String() + "/api" + r.URL.Path
	}
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
	http.HandleFunc("/", proxyHandler)

	log.Println("Proxy server is running or port 9999....")
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
