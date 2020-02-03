package svrkit

import (
	"time"
)

//TTLCache 带生存期的 cache
type TTLCache struct {
	*KVCache
}

type ttlCacheItem struct {
	data   interface{}
	expire time.Time
}

//NewTTLCache 创建一个新的 TTLCache
func NewTTLCache() *TTLCache {
	return &TTLCache{
		NewKVCache(),
	}
}

//SetWithExpire 写入
func (kv *TTLCache) SetWithExpire(k, v interface{}, expire time.Time) {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	kv.data[k] = &ttlCacheItem{
		data:   v,
		expire: expire,
	}
}

//AddIntWithExpire 设置 Int 值
func (kv *TTLCache) AddIntWithExpire(k interface{}, incr int, expire time.Time) {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	var oldVal int
	item := kv.data[k]
	if ttlItem, ok := item.(*ttlCacheItem); ok {
		if !ttlItem.expire.Before(time.Now()) {
			if v, ok := ttlItem.data.(int); ok {
				oldVal = v
			}
		}
	} else if v, ok := item.(int); ok {
		oldVal = v
	}

	kv.data[k] = &ttlCacheItem{
		data:   oldVal + incr,
		expire: expire,
	}
}

//Get 获取
func (kv *TTLCache) Get(k interface{}) interface{} {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	item := kv.data[k]
	if ttlItem, ok := item.(*ttlCacheItem); ok {
		if ttlItem.expire.Before(time.Now()) {
			delete(kv.data, k)
			return nil
		}
		return ttlItem.data

	}
	return item
}

//GetString 获取字符串类型
func (kv *TTLCache) GetString(k interface{}) string {
	item, _ := kv.Get(k).(string)
	return item
}

//GetInt 获取 int 类型
func (kv *TTLCache) GetInt(k interface{}) int {
	item, _ := kv.Get(k).(int)
	return item
}
