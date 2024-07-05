package tool

import "strings"

func EnforceHttp(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDoMain(url string) bool {
	if url == "localhost:9900" {
		return false
	}
	newUrl := strings.Replace(url, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "https://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0]
	if newUrl == "localhost:9900" {
		return false
	}
	return true
}
