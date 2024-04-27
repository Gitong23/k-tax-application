package tax

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"
)

func OpenFormFile(files []*multipart.FileHeader) ([]multipart.File, error) {
	var src []multipart.File
	for _, file := range files {
		s, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer s.Close()
		src = append(src, s)
	}
	return src, nil
}

func isCorrectHeader(record []string) bool {
	if record[0] != "totalIncome" || record[1] != "wht" || record[2] != "donation" {
		return false
	}
	return true
}

func csvTaxReq(record []string) (*TaxRequest, error) {
	if len(record) != 3 {
		return nil, fmt.Errorf("Invalid CSV file content")
	}

	income, err := strconv.ParseFloat(record[0], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid TotalIncome value")
	}

	wht, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid WHT value")
	}

	donationAmount, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid Donation value")
	}

	return &TaxRequest{
		TotalIncome: income,
		WHT:         wht,
		Allowances: []AllowanceReq{
			{
				AllowanceType: "donation",
				Amount:        donationAmount,
			},
		},
	}, nil
}

func appendTaxReq(t *[]TaxRequest, rec [][]string) error {

	for idx, r := range rec {
		if idx == 0 {
			if !isCorrectHeader(r) {
				return fmt.Errorf("Invalid CSV header")
			}
			continue
		}

		taxReq, err := csvTaxReq(r)
		if err != nil {
			return err
		}
		*t = append(*t, *taxReq)
	}

	return nil
}

func readFileCsv(f multipart.File) (records [][]string, err error) {
	reader := csv.NewReader(f)
	records, err = reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Invalid CSV file")
	}
	return records, nil
}

func fileTaxReq(src []multipart.File) ([]TaxRequest, error) {
	var taxsReq []TaxRequest
	for _, s := range src {
		records, err := readFileCsv(s)
		if err != nil {
			return nil, err
		}

		err = appendTaxReq(&taxsReq, records)
		if err != nil {
			return nil, err
		}
	}
	return taxsReq, nil
}

func NewTaxUploadResponse(t []TaxRequest, d *Deductor) *TaxUploadResponse {

	var ts []TaxUpload
	for _, tr := range t {
		i := tr.TotalIncome - d.total(tr.Allowances)
		taxUp := NewTaxUpload(tr, i)
		ts = append(ts, taxUp)
	}

	return &TaxUploadResponse{Taxs: ts}
}
