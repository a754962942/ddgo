package utils

import (
	"reflect"
	"strings"
	"unicode"
	"unsafe"
)

// SubStrAndTrim
//
//	@Description: str 从step开始截取到最后，然后去掉左边的"/"
//	@param str
//	@param step
//	@return string
func SubStrAndTrim(str, step string) string {
	index := strings.Index(str, step)
	if index < 0 {
		return ""
	}
	s := str[index+len(step):]
	left := strings.TrimLeft(s, "/")
	return left
}

// isASCII
//
//	@Description: 判断字符串是否完全为ascii码
//	@param s
//	@return bool
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// StringToBytes
//
//	@Description: string转[]byte
//	@param s
//	@return []byte
func StringToBytes2(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func StringToBytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	str := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&str))
}
