package array

import (
	"reflect"
	"strings"
)

// 查找字符是否在数组中
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// 查找数组是否有任意一个在数组中
func ArrayHasInArray(source []string, search []string) bool {
	for _, searchItem := range search {
		if InArray(searchItem, source) {
			return true
		}
	}
	return false
}

// 增加&and处理
func ArrayHasInArrayAnd(source []string, search []string) bool {
	for _, searchItem := range search {
		searchArr := strings.Split(searchItem, "&")
		if len(searchArr) > 1 {
			// 需要搜索多个
			cc := 0
			for _, searchVal := range searchArr {
				if InArray(searchVal, source) {
					cc++
				}
			}
			if cc == len(searchArr) {
				return true
			}
		} else {
			if InArray(searchItem, source) {
				return true
			}
		}

	}
	return false
}

func UniqueString(obj []string) []string {
	m := make(map[string]int)

	for _, v := range obj {
		m[v] = 1
	}

	return GetMapKeys(m)

}

// 过滤空数组
func FilterString(obj []string) []string {
	var m []string
	for _, v := range obj {
		if v != "" {
			m = append(m, v)
		}
	}
	return m

}

func GetMapKeys(data map[string]int) []string {
	ret := []string{}
	for k, _ := range data {
		ret = append(ret, k)

	}
	return ret

}

// 数组倒序函数
func Reverse(arr *[]string) {
	var temp string
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}
