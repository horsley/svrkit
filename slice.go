package svrkit

import (
	"fmt"
	"reflect"
	"strings"
)

//StringSliceUnique 文本数组去重 SliceUnique(src).([]string)
func StringSliceUnique(src []string) []string {
	return SliceUnique(src).([]string)
}

//IntSliceUnique int 数组去重 SliceUnique(src).([]int)
func IntSliceUnique(src []int) []int {
	return SliceUnique(src).([]int)
}

//SliceUnique 数组去重 自己用断言转类型
func SliceUnique(slice interface{}) interface{} {
	t := reflect.TypeOf(slice)
	v := reflect.ValueOf(slice)
	if t.Kind() != reflect.Slice {
		return slice
	}

	//start play with fire
	elementT := reflect.TypeOf(slice).Elem()
	uniqMap := reflect.MakeMap(reflect.MapOf(elementT, reflect.TypeOf(true)))
	rstSlice := reflect.MakeSlice(reflect.SliceOf(elementT), 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		if !uniqMap.MapIndex(v.Index(i)).IsValid() {
			rstSlice = reflect.Append(rstSlice, v.Index(i))
			uniqMap.SetMapIndex(v.Index(i), reflect.ValueOf(true))
		}
	}
	return rstSlice.Interface()
}

//SliceToList slice转为列表字符串，默认是逗号分割列表
func SliceToList(slice interface{}, seprator string) string {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return fmt.Sprint(slice)
	}
	if seprator == "" {
		seprator = ","
	}

	v := reflect.ValueOf(slice)
	items := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		items[i] = fmt.Sprint(v.Index(i).Interface())
	}
	return strings.Join(items, seprator)
}

//InSliceChecker 返回一个用于检查元素是否在 slice 中的闭包函数
func InSliceChecker(slice interface{}) func(needle interface{}) bool {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return func(interface{}) bool { return false }
	}

	v := reflect.ValueOf(slice)

	chkMap := make(map[interface{}]bool)
	for i := 0; i < v.Len(); i++ {
		chkMap[v.Index(i).Interface()] = true
	}
	return func(needle interface{}) bool {
		return chkMap[needle]
	}
}
