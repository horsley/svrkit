package svrkit

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

//RandNum 返回指定位数的随机数字符串
func RandNum(length int) string {
	base := 0
	if length > 1 {
		base = int(math.Pow10(length - 1))
	}
	return fmt.Sprint(rand.Intn(int(math.Pow10(length))-1-base) + base)
}
