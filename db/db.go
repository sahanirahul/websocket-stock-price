package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sensibull/stocks-api/bootconfig"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"

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
	fmt.Println("trying to connect to redis localhost address:", credentials.Url[0])
	pong := ""
	pong, err = cli.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err)
		redis_docker_addr := os.Getenv("REDIS_DOCKER_ADDR")
		if len(redis_docker_addr) == 0 {
			redis_docker_addr = "my-docker-redis:6379"
		}
		fmt.Println("retrying with redis docker address:", redis_docker_addr)
		cli = redis.NewClient(&redis.Options{
			Addr:     redis_docker_addr,
			Password: "",
			DB:       0,
		})
		pong, err = cli.Ping(context.Background()).Result()
		if err != nil {
			fmt.Println()
			log.Fatalf("redis error : %s", err.Error())
		}
	}
	fmt.Println(pong)
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
