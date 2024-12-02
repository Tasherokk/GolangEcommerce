package models

import "time"

type CacheEntry struct {
    CacheKey       string    `json:"cache_key"`
    CacheValue     string    `json:"cache_value"`
    ExpirationTime time.Time `json:"expiration_time"`
}
