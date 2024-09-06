package main

import (
	"log"
	"net/http"

	"com.zhouhc.study/client"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("this is hello ")

	w.Write([]byte("{\"code\": 1, \"data\": \"hello\"}"))

}

func main() {
	client.InitializeGProxy()

	http.HandleFunc("/hello", helloHandler)
	err := http.ListenAndServe(":9912", nil)
	if err!= nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
