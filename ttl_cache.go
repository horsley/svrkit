package svrkit

import (
	"time"
)

//TTLCache 带生存期的 cache，默认仅在读取时检查是否过期，如果内存敏感需要定期清除过期条目，请用ClearExpireItems方法
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
	kv.Set(k, &ttlCacheItem{
		data:   v,
		expire: expire,
	})
}

//AddIntWithExpire 设置 Int 值，如已有值，则加增量，重设超时时间
func (kv *TTLCache) AddIntWithExpire(k interface{}, incr int, expire time.Time) {
	var oldVal int
	if item := kv.Get(k); item != nil {
		if v, ok := item.(int); ok {
			oldVal = v
		}
	}

	kv.Set(k, &ttlCacheItem{
		data:   oldVal + incr,
		expire: expire,
	})
}

//Get 获取，如果条目已过期则返回 nil
func (kv *TTLCache) Get(k interface{}) interface{} {
	item := kv.KVCache.Get(k)
	if ttlItem, ok := item.(*ttlCacheItem); ok {
		if ttlItem.expire.Before(time.Now()) {
			kv.KVCache.Del(k)
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

//ClearExpireItems 清除过期条目，如果传入 interval 不为0，则本方法内部循环清除，不退出，你需要用 goroutine 调用。
//如果interval 为0，仅作一次遍历清理就返回
func (kv *TTLCache) ClearExpireItems(interval time.Duration) {
	for {
		keys := kv.Keys()

		for _, k := range keys {
			kv.Get(k)
		}

		if interval == 0 {
			break
		}

		time.Sleep(interval)
	}
}
