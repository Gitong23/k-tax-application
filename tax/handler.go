package tax

import (
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface {
	PersonalAllowance() (float64, error)
	DonationAllowance() (*Allowances, error)
	KreceiptAllowance() (*Allowances, error)
}

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
	result := 0.0
	for _, s := range steps {
		result += s.taxStep(netIncome)
		netIncome -= s.Max - s.Min
	}
	return result
}

func taxLevel(netIncome float64) []TaxLevel {
	var taxLevels []TaxLevel
	for _, s := range steps {
		taxLevels = append(taxLevels, TaxLevel{
			Level: fmt.Sprintf("%.0f - %.0f", s.Min, s.Max),
			Tax:   s.taxStep(netIncome),
		})
		netIncome -= s.Max - s.Min
	}
	return taxLevels
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
	initPersonalAllowance, err := h.store.PersonalAllowance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}
	incomeTax := reqTax.TotalIncome - initPersonalAllowance

	allowances := reqTax.Allowances
	for _, allowance := range allowances {

		switch allowance.AllowanceType {
		case "donation":
			donationAllowance, err := h.store.DonationAllowance()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
			}

			if allowance.Amount < donationAllowance.MinAmount {
				return c.JSON(http.StatusBadRequest, Err{Message: "Invalid donation amount"})
			}

			if allowance.Amount > donationAllowance.MaxAmount {
				incomeTax -= donationAllowance.MaxAmount
				break
			}

			incomeTax -= allowance.Amount
		case "k-receipt":
			kReceiptAllowance, err := h.store.KreceiptAllowance()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
			}

			if allowance.Amount < kReceiptAllowance.MinAmount {
				return c.JSON(http.StatusBadRequest, Err{Message: "Invalid k-receipt amount"})
			}

			if allowance.Amount > kReceiptAllowance.MaxAmount {
				incomeTax -= kReceiptAllowance.MaxAmount
				break
			}

			incomeTax -= allowance.Amount
		default:
			incomeTax -= 0
		}
	}

	wht := reqTax.WHT
	tax := calTax(incomeTax)

	var taxLevels []TaxLevel
	taxLevels = taxLevel(incomeTax)

	if wht > tax {
		return c.JSON(http.StatusOK, &TaxResponse{
			Tax:       0,
			TaxRefund: wht - tax,
			TaxLevels: taxLevels,
		})
	}

	return c.JSON(http.StatusOK, &TaxResponse{Tax: tax - wht, TaxLevels: taxLevels})
}
