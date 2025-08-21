package cache

import (
    "sync"
    "time"

)

type CacheItem struct {
    Value      interface{}
    Expiration int64
}

type Cache struct {
    items map[string]CacheItem
    mu    sync.RWMutex
}

var c = &Cache{
    items: make(map[string]CacheItem),
}

func Set(key string, value interface{}, duration time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = CacheItem{
        Value:      value,
        Expiration: time.Now().Add(duration).UnixNano(),
    }
}

func Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, found := c.items[key]
    if !found || time.Now().UnixNano() > item.Expiration {
        return nil, false
    }
    return item.Value, true
}

func Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}
