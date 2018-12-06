package strcase

import (
	"regexp"
	"strings"
)

var wordAfterDelimiter = regexp.MustCompile("[_-]([\\w])")
var delimiter = regexp.MustCompile("(.*?)[_-]([\\w])")

func ToLowerCamel(str string) (camel string) {
	if str == "" {
		return
	}
	if r := rune(str[0]); r >= 'A' && r <= 'Z' {
		str = strings.ToLower(string(r)) + str[1:]
	}
	return ToCamel(str)
}

func ToCamel(str string) (camel string) {
	if str == "" {
		return
	}
	camel = wordAfterDelimiter.ReplaceAllStringFunc(str, strings.ToUpper)
	camel = delimiter.ReplaceAllString(camel, "$1$2")
	return
}
