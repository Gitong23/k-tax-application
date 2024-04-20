package tax

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

type StepTax struct {
	Min  float64
	Max  float64
	Rate float64
}

func NewHandler() *Handler {
	return &Handler{}
}

func calTax(netIncome float64, st []StepTax) float64 {
	result := 0.0
	for _, s := range st {
		if netIncome > (s.Max - s.Min) {
			result += (s.Max - s.Min) * s.Rate
			netIncome -= (s.Max - s.Min)
			continue
		}
		result += netIncome * s.Rate
		break
	}
	return result
}

func newTax(netIncome float64, wht float64) *Tax {

	return &Tax{
		Tax: calTax(netIncome, []StepTax{
			{0, 150000, 0},
			{150000, 500000, 0.1},
			{500000, 1000000, 0.15},
			{1000000, 2000000, 0.20},
			{2000000, math.MaxFloat64, 0.35},
		}) - wht,
	}
}

func (h *Handler) CalTax(c echo.Context) error {

	taxRequest := TaxRequest{}
	if err := c.Bind(&taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//TODO:calculate income tax
	incomeTax := taxRequest.TotalIncome - 60000
	wht := taxRequest.WHT
	return c.JSON(http.StatusOK, newTax(incomeTax, wht))
}
