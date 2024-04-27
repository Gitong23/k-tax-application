package helper

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
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

func IsFilesExt(n string, f []*multipart.FileHeader) bool {
	for _, file := range f {
		if filepath.Ext(file.Filename) != n {
			return false
		}
	}
	return true
}
