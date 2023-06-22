package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sensibull/stocks-api/bootconfig"
	"strings"
	"sync"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

var cli *redis.Client
var once sync.Once

type credentials struct {
	Url      []string `json:"url"`
	Password string   `json:"password"`
	Master   string   `json:"master"`
}

func Init() {
	GetRedisClient()
}

func initRedisClient() {
	credentials, err := loadRedisConfig()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	cli = redis.NewClient(&redis.Options{
		Addr:     credentials.Url[0],
		Password: "",
		DB:       0,
	})

	pong, err := cli.Ping().Result()
	fmt.Println(pong, err)
}

func GetRedisClient() *redis.Client {
	if cli == nil {
		once.Do(initRedisClient)
	}
	return cli
}

func loadRedisConfig() (credentials, error) {
	redis := credentials{}
	byteVal, err := bootconfig.ConfigManager.Get("redis")
	if err != nil || byteVal == nil {
		if err == nil {
			err = errors.New("null_secretValue")
		}
		log.Println("error in loading aws configs:", err)
		return redis, err
	}
	var redisConfigs map[string]string
	err = json.Unmarshal(byteVal, &redisConfigs)
	if err != nil {
		log.Println("error in loading configs:", err)
		panic(err)
	}
	url, ok := redisConfigs["url"]
	if !ok {
		err := errors.New("no_url_found_in_config")
		return redis, err
	}
	redis.Url = strings.Split(url, ",")
	redis.Password = redisConfigs["password"]
	redis.Master = redisConfigs["master"]
	return redis, nil
}
