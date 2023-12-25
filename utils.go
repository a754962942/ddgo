package ddgo

import "strings"

// SubStrAndTrim
//
//	@Description:
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
