package svrkit

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

//PrettySize 美化尺寸信息显示
func PrettySize(b int64) string {
	if b < 1024 {
		return fmt.Sprint(b, "B")
	} else if b < 1024*1024 {
		return fmt.Sprintf("%.2gKB", float64(b)/float64(1024))
	} else if b < 1024*1024*1024 {
		return fmt.Sprintf("%.2gMB", float64(b)/float64(1024*1024))
	}
	return fmt.Sprintf("%.2gGB", float64(b)/float64(1024*1024*1024))
}

//PrettyDate 优化日期显示
func PrettyDate(t time.Time) string {
	y, m, d := t.Date()
	ny, nm, nd := time.Now().Date()
	if y == ny && m == nm && d == nd {
		return t.Format("今天 15:04")
	} else if y == ny {
		return t.Format("01/02 15:04")
	} else {
		return t.Format("2006/01/02 15:04")
	}
}

//Base64 Base64转码
func Base64(src string) string {
	return base64.URLEncoding.EncodeToString([]byte(src))
}

//MustInt 强转 int
func MustInt(src string) int {
	i, _ := strconv.Atoi(src)
	return i
}
