package tax

import "mime/multipart"

type Service interface {
	OpenCsvFile(filePath string) ([][]string, error)
}

type CsvReader struct {
	file []*multipart.FileHeader
}
