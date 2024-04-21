package tax

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface{}

type StepTax struct {
	Min  float64
	Max  float64
	Rate float64
}

type Err struct {
	Message string `json:"message"`
}

func NewHandler(db Storer) *Handler {
	return &Handler{store: db}
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

func calTax(netIncome float64) float64 {

	steps := []StepTax{
		{0, 150000, 0},
		{150000, 500000, 0.1},
		{500000, 1000000, 0.15},
		{1000000, 2000000, 0.20},
		{2000000, math.MaxFloat64, 0.35},
	}

	result := 0.0
	for _, s := range steps {
		result += s.taxStep(netIncome)
		netIncome -= s.Max - s.Min
	}
	return result
}

func newTax(netIncome float64, wht float64) *Tax {
	return &Tax{Tax: calTax(netIncome) - wht}
}

func (h *Handler) CalTax(c echo.Context) error {

	reqTax := TaxRequest{}
	if err := c.Bind(&reqTax); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if reqTax.WHT > reqTax.TotalIncome || reqTax.WHT < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid WHT value"})
	}

	//TODO:calculate income tax
	incomeTax := reqTax.TotalIncome - 60000

	allowances := reqTax.Allowances
	for _, allowance := range allowances {

		switch allowance.AllowanceType {
		case "donation":
			if allowance.Amount > 100000 {
				incomeTax -= 100000
			} else {
				incomeTax -= allowance.Amount
			}
		case "k-receipt":
			incomeTax -= 50000 //initial value
		default:
			incomeTax -= 0
		}
	}

	wht := reqTax.WHT
	return c.JSON(http.StatusOK, newTax(incomeTax, wht))
}
