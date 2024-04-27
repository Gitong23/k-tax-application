package tax

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"
)

type Csv struct {
	TaxReq *[]TaxRequest
}

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

func convCsv(src []multipart.File) ([]TaxRequest, error) {
	var taxsReq []TaxRequest
	for _, s := range src {
		reader := csv.NewReader(s)
		records, err := reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("Invalid CSV file")
		}

		for idx, record := range records {
			//check header file
			if idx == 0 {
				if record[0] != "totalIncome" || record[1] != "wht" || record[2] != "donation" {
					return nil, fmt.Errorf("Invalid CSV header")
				}
				continue
			}

			//check content file
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

			taxReq := TaxRequest{
				TotalIncome: income,
				WHT:         wht,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        donationAmount,
					},
				},
			}
			taxsReq = append(taxsReq, taxReq)
		}
	}
	return taxsReq, nil
}
