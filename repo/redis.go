package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func RedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "docker.for.mac.localhost:6379",
	})

	check, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	if check != "PONG" {
		return nil, fmt.Errorf("redis server not responding")
	}
	fmt.Println("Ping: ", check)

	return rdb, nil
}

func ReadDbWithCache(readDbFunc func() ([]byte, error), key string, client *redis.Client) ([]byte, error) {
	redisData, err := client.Get(ctx, key).Result()
	ifReadFromDb := false

	if err == redis.Nil || err != nil {
		ifReadFromDb = true
		// fmt.Println(key, " does not exist")
	}
	// fmt.Println(key, "-- redisdata ---", redisData)

	if ifReadFromDb {
		setRedis, err1 := readDbFunc()
		if err1 != nil {
			return nil, err1
		}
		_, err = client.Set(ctx, key, setRedis, 100*time.Second).Result()
		if err != nil {
			fmt.Println(err)
		}
		return setRedis, err
	}
	return []byte(redisData), nil
}
