package tax

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, Err{Message: "Invalid request body"})
	}
	return nil
}

func (t *TaxRequest) validatWht() error {
	if t.WHT > t.TotalIncome || t.WHT < 0 {
		return fmt.Errorf("Invalid WHT value")
	}
	return nil
}

func checkMultiWht(t []TaxRequest) error {
	for _, taxReq := range t {
		err := taxReq.validatWht()
		if err != nil {
			return err
		}
	}
	return nil
}
