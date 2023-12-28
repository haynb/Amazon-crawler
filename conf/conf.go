package conf

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	VerifyUserKey     string `yaml:"verifyUserKey"`
	MongoDBAddr       string `yaml:"mongoDBAddr"`
	MongoDBUser       string `yaml:"mongoDBUser"`
	MongoDBPwd        string `yaml:"mongoDBPwd"`
	MongoDBDatabase   string `yaml:"mongoDBDatabase"`
	MongoDBCollection string `yaml:"mongoDBCollection"`
}

var Conf Config

func init() {
	file, err := os.Open("./conf.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = yaml.NewDecoder(file).Decode(&Conf)
	if err != nil {
		panic(err)
	}
}
