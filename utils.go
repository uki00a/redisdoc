package main

import (
	"regexp"
	"strings"
)

func RemoveConsecutiveSpaces(s string) string {
	return consecutiveSpacesRe.ReplaceAllString(s, " ")
}

func StartsWithPunctuation(s string) bool {
	return strings.HasPrefix(s, ".") ||
		strings.HasPrefix(s, ",") ||
		strings.HasPrefix(s, ":")
}

var consecutiveSpacesRe *regexp.Regexp

func init() {
	consecutiveSpacesRe = regexp.MustCompile(`\s{2,}`)
}
