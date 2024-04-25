package tax

import (
	"fmt"
	"math"

	"github.com/Gitong23/assessment-tax/helper"
)

type StepTax struct {
	Min  float64
	Max  float64
	Rate float64
}

var steps = []StepTax{
	{0, 150000, 0},
	{150000, 500000, 0.1},
	{500000, 1000000, 0.15},
	{1000000, 2000000, 0.20},
	{2000000, math.MaxFloat64, 0.35},
}

func taxLevel(netIncome float64) []TaxLevel {
	var taxLevels []TaxLevel
	for idx, s := range steps {

		var level string
		if idx == 0 {
			level = fmt.Sprintf("0 - %s", helper.Comma(s.Max))
		}

		if idx == len(steps)-1 {
			level = fmt.Sprintf("%s ขึ้นไป", helper.Comma(s.Min))
		}

		if idx > 0 && idx < len(steps)-1 {
			level = fmt.Sprintf("%s - %s", helper.Comma(s.Min+1), helper.Comma(s.Max))
		}

		taxLevels = append(taxLevels, TaxLevel{
			Level: level,
			Tax:   s.taxStep(netIncome),
		})
		netIncome -= s.Max - s.Min
	}
	return taxLevels
}

func calLevelTax(netIncome float64) float64 {
	result := 0.0
	for _, s := range steps {
		result += s.taxStep(netIncome)
		netIncome -= s.Max - s.Min
	}
	return result
}

func (s *StepTax) taxStep(amount float64) float64 {
	if amount <= 0 {
		return 0
	}

	if amount > (s.Max - s.Min) {
		return (s.Max - s.Min) * s.Rate
	}

	return amount * s.Rate
}

func NewTaxResponse(wht float64, income float64) TaxResponse {
	tax := calLevelTax(income)

	var taxLevels []TaxLevel
	taxLevels = taxLevel(income)

	if wht > tax {
		return TaxResponse{
			Tax:       0,
			TaxRefund: wht - tax,
			TaxLevels: taxLevels,
		}
	}

	return TaxResponse{
		Tax:       tax - wht,
		TaxLevels: taxLevels,
	}
}

func NewTaxUpload(taxReq TaxRequest, income float64) TaxUpload {
	tax := calLevelTax(income)

	if taxReq.WHT > tax {
		var refund float64
		refund = taxReq.WHT - tax
		return TaxUpload{
			TotalIncome: taxReq.TotalIncome,
			Tax:         0,
			TaxRefund:   &refund,
		}
	}

	return TaxUpload{
		TotalIncome: taxReq.TotalIncome,
		Tax:         tax - taxReq.WHT,
		TaxRefund:   nil,
	}
}
