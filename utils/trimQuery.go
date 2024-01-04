package utils

import (
	"regexp"
	"strings"
)

func TrimQuery(str string) string {
	compile := regexp.MustCompile("\\?.*")
	matchString := compile.MatchString(str)
	if !matchString {
		return str
	}
	findString := compile.FindString(str)
	return strings.Replace(str, findString, "", -1)
}
