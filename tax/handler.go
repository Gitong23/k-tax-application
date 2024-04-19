package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CalTax(c echo.Context) error {

	// var taxRequest TaxRequest
	taxRequest := TaxRequest{}
	if err := c.Bind(&taxRequest); err != nil {
		return err
	}

	if taxRequest.TotalIncome < 150000 {
		return c.JSON(http.StatusOK, Tax{Tax: 0})
	}

	// tax := Tax{}
	// tax.Tax = taxRequest.TotalIncome - taxRequest.WHT
	// for _, allowance := range taxRequest.Allowances {
	// 	tax.Tax -= allowance.Amount
	// }

	return c.JSON(http.StatusOK, Tax{Tax: 29000})
}
