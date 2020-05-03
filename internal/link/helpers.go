package link

import "strings"

// Parse parses the url and returns the key and param.
// key: the key of the link record without leading and trailing "/".
// param: the param passed in the url
// "/abc.efg/yoyoyo" -> ("abc.efg", "yoyoyo")
func Parse(url string) (key string, param string) {
	// remove leading and trailing `sep`
	url = strings.Trim(url, "/")
	return Pop(url, "/")
}

// Pop returns the prefix before the first `sep` and the remaining string.
func Pop(str string, sep string) (left string, remain string) {
	// split from the first `sep`
	idx := strings.Index(str, sep)
	if idx == -1 {
		left = str
	} else {
		left = str[:idx]
		remain = str[idx+len(sep):]
	}
	return
}
