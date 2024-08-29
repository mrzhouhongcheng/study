package main

import (
	"net/http"
)

func main() {

	fileServer := http.FileServer(http.Dir("./"))

	http.Handle("/", fileServer)
	if err := http.ListenAndServe(":8888", nil); err != nil {
		panic(err)
	}
}
