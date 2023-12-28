package database

import (
	"amson/conf"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var Client *mongo.Client // 包级别的 Client 变量

func init() {
	// 读取配置文件
	addr := conf.Conf.MongoDBAddr
	user := conf.Conf.MongoDBUser
	pwd := conf.Conf.MongoDBPwd
	// 注意这里的改动：去掉了 := 来确保我们赋值到包级别的Client变量
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", user, pwd, addr))
	// 连接到MongoDB
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
}
func Disconnect() {
	Client.Disconnect(context.TODO())
	Client = nil
	log.Println("Disconnected from MongoDB!")
	return
}
func GetCollection() *mongo.Collection {
	return Client.Database(conf.Conf.MongoDBDatabase).Collection(conf.Conf.MongoDBCollection)
}
func InsertOne(data interface{}) {
	collection := GetCollection()
	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.InsertedID)
}
