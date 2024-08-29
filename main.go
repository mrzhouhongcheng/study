package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const url = "mongodb://10.88.19.100:27017,10.88.19.101:27017,10.88.19.102:27017/fish?replicaSet=rs01"
const database_name = "zhouhc_test"

func main() {

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to MongoDB!")

	args := os.Args
	fileName := args[1]
	tableName := strings.Replace(fileName, ".js", "", 1)
	fmt.Println("table name is :", tableName)
	var jsonStr = load(fileName)
	var jsonData []string

	json.Unmarshal([]byte(jsonStr), &jsonData)

	for _, item := range jsonData {
		var insertData bson.M
		bson.UnmarshalExtJSON([]byte(item), true, &insertData)
		fmt.Println(insertData)
		insertToDB(insertData, tableName, client)
	}
}

func insertToDB(data bson.M, tableName string, client *mongo.Client) {
	collection := client.Database(database_name).Collection(tableName)

	count, err := collection.CountDocuments(context.TODO(), bson.M{"_id": data["_id"]})
	if count > 0 {
		collection.DeleteOne(context.TODO(), bson.M{"_id": data["_id"]})
		fmt.Println("删除一条数据")
	}
	collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", data)
}

func load(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("文件读取失败: ", path)
	}
	return string(data)
}
