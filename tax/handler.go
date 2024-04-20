package tax

import (
	"fmt"
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
	{150000, 500000, 0.1},
	{500000, 1000000, 0.15},
	{1000000, 2000000, 0.20},
	{2000000, math.MaxFloat64, 0.35},
}

func (h *Handler) CalTax(c echo.Context) error {
	// var taxRequest TaxRequest
	taxRequest := TaxRequest{}
	if err := c.Bind(&taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//TODO:calculate income tax
	incomeTax := taxRequest.TotalIncome - 60000

	fmt.Printf("incomeTax: %.2f\n", incomeTax)
	fmt.Printf("$$$$$$$$$$$$$$$$\n")
	//TODO:calculate tax
	sumTax := 0.0
	for idx, s := range step {

		if incomeTax > (s.Max - s.Min) {
			sumTax += (s.Max - s.Min) * s.Rate
			incomeTax -= (s.Max - s.Min)
			fmt.Printf("range: %.2f - %.2f\n", s.Min, s.Max)
			fmt.Printf("Amount Tax: %.2f\n", (s.Max-s.Min)*s.Rate)
			fmt.Printf("sumTax: %.2f\n", sumTax)
			fmt.Printf("incomeTax: %.2f\n", incomeTax)
			fmt.Printf("---------%d-----------\n", idx+1)

			continue
		}
		sumTax += incomeTax * s.Rate
		fmt.Printf("sumTax: %v\n", sumTax)
		break
	}

	return c.JSON(http.StatusOK, Tax{Tax: sumTax})
}
