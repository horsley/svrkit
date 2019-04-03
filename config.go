package svrkit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

//ConfigEnvSuffix 开发配置覆盖后缀
var ConfigEnvSuffix = "dev"

//ConfigBoolTrueValue 布尔值开关的真值字符串
var ConfigBoolTrueValue = "yes"

//Config 配置信息 带读写锁哦
type Config struct {
	kv    map[string][]string
	lock  sync.RWMutex
	isDev bool
}

//LoadFromReader 从 reader 读取
func (c *Config) LoadFromReader(f io.Reader) {
	newKv := c.readKV(f)

	c.lock.Lock()
	defer c.lock.Unlock()
	if dev, ok := newKv[ConfigEnvSuffix]; ok && len(dev) > 0 && dev[0] == ConfigBoolTrueValue {
		c.isDev = true
	}

	c.kv = newKv
}

func (c *Config) readKV(f io.Reader) map[string][]string {
	kvRe := regexp.MustCompile(`(.*?)=(.*?)$`)
	newKv := make(map[string][]string)

	s := bufio.NewScanner(f)
	for s.Scan() {
		if kvRe.MatchString(s.Text()) {
			m := kvRe.FindStringSubmatch(s.Text())
			k := strings.TrimSpace(m[1])
			v := strings.TrimSpace(m[2])

			if old, ok := newKv[k]; ok {
				newKv[k] = append(old, v)
			} else {
				newKv[k] = []string{v}
			}

		}
	}
	return newKv
}

//Load 从文件读入 如果有filename.{ConfigEnvSuffix}文件存在，其中的配置会覆盖进来
func (c *Config) Load(filename string) error {

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	c.LoadFromReader(f)

	for {
		devOverride := filename + "." + ConfigEnvSuffix
		if _, err := os.Stat(devOverride); err == nil {
			o, err := os.Open(devOverride)
			if err != nil {
				break
			}
			defer o.Close()

			overrideKV := c.readKV(o)
			for k, v := range overrideKV {
				c.kv[k] = v
			}
		}
		break
	}

	return nil
}

//LoadFromString 从字符串解析
func (c *Config) LoadFromString(conf string) {
	f := strings.NewReader(conf)
	c.LoadFromReader(f)
}

//GetArray 获取文本数组
func (c *Config) GetArray(key string) []string {
	return c.GetWithDefault(key, []string{})
}

//GetWithDefault 读取文本数组 带默认项
func (c *Config) GetWithDefault(key string, def []string) []string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if v, ok := c.kv[key+"."+ConfigEnvSuffix]; ok && c.isDev { //存在开发覆盖
		return v
	}

	if v, ok := c.kv[key]; ok {
		return v
	}
	return def
}

//GetFirst 取首个出现的值
func (c *Config) GetFirst(key string) string {
	if array := c.GetArray(key); len(array) > 0 {
		return array[0]
	}
	return ""
}

//Get 取首个出现的值 GetFirst 的别名
func (c *Config) Get(key string) string {
	return c.GetFirst(key)
}

//GetFirstWithDefault 取首个出现的值 带默认配置
func (c *Config) GetFirstWithDefault(key, def string) string {
	return c.GetWithDefault(key, []string{def})[0]
}

//GetInt 取整数值配置
func (c *Config) GetInt(key string) int {
	return c.mustInt(c.GetFirst(key))
}

//GetFloat 取浮点值配置
func (c *Config) GetFloat(key string) float64 {
	f, err := strconv.ParseFloat(c.Get(key), 64)
	if err != nil {
		return 0
	}
	return f
}

//GetIntWithDefault 取整数值配置带默认值
func (c *Config) GetIntWithDefault(key string, def int) int {
	s := c.GetWithDefault(key, []string{fmt.Sprint(def)})
	return c.mustInt(s[0])
}

func (c *Config) mustInt(src string) int {
	r, _ := strconv.Atoi(src)
	return r
}

//GetBool 取布尔值配置
func (c *Config) GetBool(key string) bool {
	return c.GetBoolWithDefault(key, false)
}

//GetBoolWithDefault 带默认值取布尔值配置
func (c *Config) GetBoolWithDefault(key string, def bool) bool {
	var defStr string
	if def {
		defStr = ConfigBoolTrueValue
	}

	return strings.ToLower(c.GetWithDefault(key, []string{defStr})[0]) == ConfigBoolTrueValue
}
