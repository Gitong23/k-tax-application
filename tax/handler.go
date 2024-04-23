package tax

import (
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Gitong23/assessment-tax/helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface {
	PersonalAllowance() (*Allowances, error)
	DonationAllowance() (*Allowances, error)
	KreceiptAllowance() (*Allowances, error)
	UpdateInitPersonalAllowance(amount float64) error
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

func (h *Handler) CalTax(c echo.Context) error {

	reqTax := TaxRequest{}
	if err := c.Bind(&reqTax); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if reqTax.WHT > reqTax.TotalIncome || reqTax.WHT < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid WHT value"})
	}

	//TODO:calculate income tax
	personal, err := h.store.PersonalAllowance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}
	incomeTax := reqTax.TotalIncome - personal.InitAmount

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
		// case "k-receipt":
		// 	kReceiptAllowance, err := h.store.KreceiptAllowance()
		// 	if err != nil {
		// 		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
		// 	}

		// 	if allowance.Amount < kReceiptAllowance.MinAmount {
		// 		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid k-receipt amount"})
		// 	}

		// 	if allowance.Amount > kReceiptAllowance.MaxAmount {
		// 		incomeTax -= kReceiptAllowance.MaxAmount
		// 		break
		// 	}

		// 	incomeTax -= allowance.Amount
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

func (h *Handler) SetPersonalDeduction(c echo.Context) error {
	reqAmount := DeductionReq{}
	if err := c.Bind(&reqAmount); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	exitPersonal, err := h.store.PersonalAllowance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	if reqAmount.Amount < exitPersonal.MinAmount || reqAmount.Amount > exitPersonal.MaxAmount {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid personal deduction amount"})
	}

	//TODO: update personal deduction
	if err := h.store.UpdateInitPersonalAllowance(reqAmount.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	newPersonal, err := h.store.PersonalAllowance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, &DeductionRes{PersonalDeduction: newPersonal.InitAmount})
}

func (h *Handler) UploadCsv(c echo.Context) error {
	// Read form data
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	var taxs []TaxRequest

	files := form.File["taxFile"]
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		if ext != ".csv" {
			return c.JSON(http.StatusBadRequest, "Only CSV files are allowed")
		}

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		//Parse CSV
		reader := csv.NewReader(src)
		records, err := reader.ReadAll()
		if err != nil {
			return err
		}

		for idx, record := range records {

			if idx == 0 {
				if record[0] != "totalIncome" || record[1] != "wht" || record[2] != "donation" {
					return c.JSON(http.StatusBadRequest, "Invalid CSV header")
				}
				continue
			}

			income, err := strconv.ParseFloat(record[0], 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, "Invalid TotalIncome value")
			}

			wht, err := strconv.ParseFloat(record[1], 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, "Invalid WHT value")
			}

			donationAmount, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, "Invalid Donation value")
			}

			tax := TaxRequest{
				TotalIncome: income,
				WHT:         wht,
				Allowances: []AllowanceReq{
					{
						AllowanceType: record[2],
						Amount:        donationAmount,
					},
				},
			}

			taxs = append(taxs, tax)
		}

	}

	//Convert CSV to Struct

	return c.JSON(http.StatusOK, taxs)

}
