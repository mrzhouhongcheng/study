package main 

import (
	"io"
	"log"
	"net/http"
)


// 服务器 


func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	targetUrl := r.Header[http.CanonicalHeaderKey("proxyURL")][0]
	log.Printf("[INFO] Request url is %s\n", targetUrl)
	req, err := http.NewRequest(r.Method, targetUrl, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		if key == "proxyURL" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	// 启动一个代理服务器； 转发所有的请求；
	http.HandleFunc("/", handleRequestAndRedirect)
	log.Println("Server is listering on port 9999....")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
