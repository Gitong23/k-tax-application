package tax

import (
	"fmt"
	"net/http"

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
	UpdateMaxAmountKreceipt(amount float64) error
}

type Err struct {
	Message string `json:"message"`
}

func NewHandler(db Storer) *Handler {
	return &Handler{store: db}
}

func (h *Handler) CalTax(c echo.Context) error {

	reqTax := TaxRequest{}
	if err := c.Bind(&reqTax); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := reqTax.validatWht()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	deductor, err := NewDeductor(h.store)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	err = deductor.checkMinAllowanceReq(reqTax.Allowances)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	incomeTax := reqTax.TotalIncome - deductor.total(reqTax.Allowances)
	return c.JSON(http.StatusOK, NewTaxResponse(reqTax.WHT, incomeTax))
}

func validateInitPersonalDeduction(s Storer, amount float64) error {

	p, err := s.PersonalAllowance()
	if err != nil {
		return fmt.Errorf("Internal Server Error")
	}

	if amount < p.MinAmount || amount > p.MaxAmount {
		return fmt.Errorf("Invalid personal deduction amount")
	}

	return nil
}

func (h *Handler) UpdateInitPersonalDeduct(c echo.Context) error {
	reqAmount := DeductionReq{}
	if err := c.Bind(&reqAmount); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := validateInitPersonalDeduction(h.store, reqAmount.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if err := h.store.UpdateInitPersonalAllowance(reqAmount.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	newPersonal, err := h.store.PersonalAllowance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, &InitPersonalDeductRes{PersonalDeduction: newPersonal.InitAmount})
}

func (h *Handler) UpdateMaxKreceiptDeduct(c echo.Context) error {

	reqAmount := DeductionReq{}
	if err := c.Bind(&reqAmount); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	exitKreceipt, err := h.store.KreceiptAllowance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	if reqAmount.Amount < exitKreceipt.MinAmount || reqAmount.Amount > exitKreceipt.LimitMaxAmount {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid K-receipt deduction amount"})
	}

	if err := h.store.UpdateMaxAmountKreceipt(reqAmount.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	newKreceipt, err := h.store.KreceiptAllowance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, &MaxKreceiptRes{Kreceipt: newKreceipt.MaxAmount})
}

func (h *Handler) UploadCsv(c echo.Context) error {
	// Read form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid form data")
	}

	files := form.File["taxFile"]
	if !helper.IsFilesExt(".csv", files) {
		return c.JSON(http.StatusBadRequest, "Only CSV files are allowed")
	}

	src, err := OpenFormFile(files)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	taxesReq, err := fileTaxReq(src)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	err = checkMultiWht(taxesReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	deductor, err := NewDeductor(h.store)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	err = deductor.checkMinMultiTaxReq(taxesReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, NewTaxUploadResponse(taxesReq, deductor))
}
