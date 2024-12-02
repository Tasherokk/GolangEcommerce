package controllers

import (
    "time"
	"database/sql"

    "github.com/Tasherokk/GolangProject.git/database"
)

func SetCacheEntry(key string, value string, duration time.Duration) error {
    expiration := time.Now().Add(duration)
    _, err := database.DB.Exec(`
        INSERT INTO Cache (cache_key, cache_value, expiration_time)
        VALUES ($1, $2, $3)
        ON CONFLICT (cache_key) DO UPDATE
        SET cache_value = $2, expiration_time = $3`,
        key, value, expiration)
    return err
}

func GetCacheEntry(key string) (string, error) {
    var value string
    var expiration time.Time
    err := database.DB.QueryRow(`
        SELECT cache_value, expiration_time FROM Cache WHERE cache_key=$1`, key).
        Scan(&value, &expiration)
    if err != nil {
        return "", err
    }

    if time.Now().After(expiration) {
        // Кеш истек
        return "", sql.ErrNoRows
    }

    return value, nil
}
