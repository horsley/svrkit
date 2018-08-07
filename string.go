package svrkit

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
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

//SubstringByFlag 通过左右标记切割文本，返回第一个匹配
func SubstringByFlag(src, left, right string) string {
	slice1 := strings.SplitN(src, left, 2)
	if len(slice1) < 2 {
		return ""
	}
	slice2 := strings.SplitN(slice1[1], right, 2)
	if len(slice1) < 2 {
		return ""
	}

	return slice2[0]

}
