package filter

import (
	"regexp"
)

func Filter(url string, regUrls []string) bool {
	if len(url) == 0 {
		return false
	}

	for _, regUrl := range regUrls {
		reg := regexp.MustCompile(regUrl)
		match := reg.MatchString(url)
		if match {
			return true
		}
	}

	return false
}
