package helper

import (
	"fmt"
	"regexp"
)

func Comma(num float64) string {
	str := fmt.Sprintf("%.0f", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != str; {
		n = str
		str = re.ReplaceAllString(str, "$1,$2")
	}
	return str
}
