package tax

import (
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
	UpdateInitPersonalAllowance(amount float64) (*Allowances, error)
	UpdateMaxAmountKreceipt(amount float64) (*Allowances, error)
}

type Err struct {
	Message string `json:"message"`
}

func NewHandler(db Storer) *Handler {
	return &Handler{store: db}
}

func (h *Handler) Tax(c echo.Context) error {

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

func (h *Handler) UpdateInitPersonalDeduct(c echo.Context) error {
	reqAmount := DeductionReq{}
	if err := c.Bind(&reqAmount); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	status, err := validateInitPersonalDeduction(h.store, reqAmount.Amount)
	if err != nil {
		return c.JSON(status, Err{Message: err.Error()})
	}
	p, err := h.store.UpdateInitPersonalAllowance(reqAmount.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, &InitPersonalDeductRes{PersonalDeduction: p.InitAmount})
}

func (h *Handler) UpdateMaxKreceiptDeduct(c echo.Context) error {

	reqAmount := DeductionReq{}
	if err := c.Bind(&reqAmount); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	status, err := validateMaxKreceipt(h.store, reqAmount.Amount)
	if err != nil {
		return c.JSON(status, Err{Message: err.Error()})
	}

	k, err := h.store.UpdateMaxAmountKreceipt(reqAmount.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, &MaxKreceiptRes{Kreceipt: k.MaxAmount})
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
