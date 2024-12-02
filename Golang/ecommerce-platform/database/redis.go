// database/redis.go
package database

import (
    "context"
    "github.com/go-redis/redis/v8"
	"log"
)

var Ctx = context.Background()
var RDB *redis.Client

func InitRedis() {
    RDB = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // Нет пароля по умолчанию
        DB:       0,  // Используем базу данных по умолчанию
    })

    _, err := RDB.Ping(Ctx).Result()
    if err != nil {
        panic(err)
    }

    log.Println("Successfully connected to Redis!")
}
