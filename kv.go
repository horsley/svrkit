package svrkit

import (
	"sync"
)

//KVCache 线程安全的 map
type KVCache struct {
	lock sync.RWMutex
	data map[interface{}]interface{}
}

//NewKVCache 创建一个新的 KVCache
func NewKVCache() *KVCache {
	return &KVCache{
		data: make(map[interface{}]interface{}),
	}
}

//Set 写入
func (kv *KVCache) Set(k, v interface{}) {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	kv.data[k] = v
}

//Del 删除
func (kv *KVCache) Del(k interface{}) {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	delete(kv.data, k)
}

//Count 获取数量
func (kv *KVCache) Count() int {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	return len(kv.data)
}

//Get 获取
func (kv *KVCache) Get(k interface{}) interface{} {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	return kv.data[k]
}

//GetString 获取字符串类型
func (kv *KVCache) GetString(k interface{}) string {
	item, _ := kv.Get(k).(string)
	return item
}

//GetInt 获取 int 类型
func (kv *KVCache) GetInt(k interface{}) int {
	item, _ := kv.Get(k).(int)
	return item
}

//Has has 方法比起 get 方法多了一个是否存在的布尔值返回，这样你可以区分到底是不存在还是存了 nil 进去
func (kv *KVCache) Has(k interface{}) (bool, interface{}) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	data, has := kv.data[k]
	return has, data
}

//Clear 清空数据
func (kv *KVCache) Clear() {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	kv.data = make(map[interface{}]interface{})
}
