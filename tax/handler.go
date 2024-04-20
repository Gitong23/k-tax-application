package tax

import (
	"fmt"
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

func NewTax(netIncome float64) *Tax {
	step := []StepTax{
		{0, 150000, 0},
		{150000, 500000, 0.1},
		{500000, 1000000, 0.15},
		{1000000, 2000000, 0.20},
		{2000000, math.MaxFloat64, 0.35},
	}

	sumTax := 0.0
	for _, s := range step {

		if netIncome > (s.Max - s.Min) {
			sumTax += (s.Max - s.Min) * s.Rate
			netIncome -= (s.Max - s.Min)
			continue
		}
		sumTax += netIncome * s.Rate
		fmt.Printf("sumTax: %v\n", sumTax)
		break
	}

	// return sumTax
	// return Tax{Tax: sumTax}
	return &Tax{Tax: sumTax}
}

func (h *Handler) CalTax(c echo.Context) error {

	taxRequest := TaxRequest{}
	if err := c.Bind(&taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//TODO:calculate income tax
	incomeTax := taxRequest.TotalIncome - 60000
	return c.JSON(http.StatusOK, NewTax(incomeTax))
}
