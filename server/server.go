package main

import (
	"fmt"
	"log"
	"net/http"

	"com.zhouhc.study/client"
	"github.com/spf13/viper"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("this is hello ")

	w.Write([]byte("{\"code\": 1, \"data\": \"hello\"}"))

}

func main() {
	defer fmt.Println("out server is listening")

	viper.SetConfigName("application-cli.yml")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	client.InitializeGProxy()

	http.HandleFunc("/hello", helloHandler)
	port := viper.GetInt("gproxy.port")
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
