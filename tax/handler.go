package tax

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type StepTax struct {
	Min  float64
	Max  float64
	Rate float64
}

var step = []StepTax{
	{0, 150000, 0},
	{150001, 500000, 0.1},
	{500001, 1000000, 0.15},
	{1000001, 2000000, 0.20},
	{2000001, math.MaxFloat64, 0.35},
}

func (h *Handler) CalTax(c echo.Context) error {
	// var taxRequest TaxRequest
	taxRequest := TaxRequest{}
	if err := c.Bind(&taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	incomeTax := taxRequest.TotalIncome - 60000

	//calculate tax
	sumTax := 0.0
	for _, s := range step {
		if incomeTax > s.Max {
			sumTax += (s.Max - s.Min) * s.Rate
			incomeTax -= s.Max
		} else {
			sumTax += incomeTax * s.Rate
			break
		}
	}

	return c.JSON(http.StatusOK, Tax{Tax: sumTax})
}
