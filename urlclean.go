package Usom

import "strings"

/*
clean up the url distorting characters
*/
func Cleanurl(url string) string {
	if strings.Contains(url, "https://") {
		return Cleanurl(strings.Split(url, "https://")[1])
	}
	if strings.Contains(url, "http://") {
		return Cleanurl(strings.Split(url, "http://")[1])
	}
	if strings.ContainsAny(url, "/") {
		return strings.Split(url, "/")[0]
	}
	return url
}
