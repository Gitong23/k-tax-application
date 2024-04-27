package tax

import (
	"encoding/csv"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

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

	if reqTax.WHT > reqTax.TotalIncome || reqTax.WHT < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid WHT value"})
	}

	deductor, err := NewDeductor(h.store)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	err = deductor.checkMinAllowanceReq(reqTax.Allowances)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	incomeTax := reqTax.TotalIncome - deductor.deductIncome(reqTax.Allowances)
	return c.JSON(http.StatusOK, NewTaxResponse(reqTax.WHT, incomeTax))
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

	return c.JSON(http.StatusOK, &InitPersonalDeductRes{PersonalDeduction: newPersonal.InitAmount})
}

func (h *Handler) SetKreceiptDeduction(c echo.Context) error {

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

	var taxsReq []TaxRequest

	files := form.File["taxFile"]
	//check file extension
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		if ext != ".csv" {
			return c.JSON(http.StatusBadRequest, "Only CSV files are allowed")
		}
	}

	//Open file form multipart
	var src []multipart.File
	for _, file := range files {
		s, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Internal Server Error")
		}
		defer s.Close()
		src = append(src, s)
	}

	//Parse CSV from src
	for _, s := range src {
		reader := csv.NewReader(s)
		records, err := reader.ReadAll()
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid CSV file")
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

			taxReq := TaxRequest{
				TotalIncome: income,
				WHT:         wht,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        donationAmount,
					},
				},
			}

			taxsReq = append(taxsReq, taxReq)
		}
	}

	// for _, file := range files {
	// 	// Open uploaded file
	// 	src, err := file.Open()
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	// 	}
	// 	defer src.Close()

	// 	//Parse CSV
	// 	reader := csv.NewReader(src)
	// 	records, err := reader.ReadAll()
	// 	if err != nil {
	// 		return c.JSON(http.StatusBadRequest, "Invalid CSV file")
	// 	}

	// 	for idx, record := range records {

	// 		if idx == 0 {
	// 			if record[0] != "totalIncome" || record[1] != "wht" || record[2] != "donation" {
	// 				return c.JSON(http.StatusBadRequest, "Invalid CSV header")
	// 			}
	// 			continue
	// 		}

	// 		income, err := strconv.ParseFloat(record[0], 64)
	// 		if err != nil {
	// 			return c.JSON(http.StatusBadRequest, "Invalid TotalIncome value")
	// 		}

	// 		wht, err := strconv.ParseFloat(record[1], 64)
	// 		if err != nil {
	// 			return c.JSON(http.StatusBadRequest, "Invalid WHT value")
	// 		}

	// 		donationAmount, err := strconv.ParseFloat(record[2], 64)
	// 		if err != nil {
	// 			return c.JSON(http.StatusBadRequest, "Invalid Donation value")
	// 		}

	// 		taxReq := TaxRequest{
	// 			TotalIncome: income,
	// 			WHT:         wht,
	// 			Allowances: []AllowanceReq{
	// 				{
	// 					AllowanceType: "donation",
	// 					Amount:        donationAmount,
	// 				},
	// 			},
	// 		}

	// 		taxsReq = append(taxsReq, taxReq)
	// 	}

	// }

	//prepare to calculation
	deductor, err := NewDeductor(h.store)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	for _, taxReq := range taxsReq {
		if taxReq.WHT > taxReq.TotalIncome || taxReq.WHT < 0 {
			return c.JSON(http.StatusBadRequest, Err{Message: "Invalid WHT value"})
		}

		err = deductor.checkMinAllowanceReq(taxReq.Allowances)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
		}
	}

	var taxsRes []TaxUpload
	for _, taxReq := range taxsReq {
		incomeTax := taxReq.TotalIncome - deductor.deductIncome(taxReq.Allowances)
		taxUpload := NewTaxUpload(taxReq, incomeTax)
		taxsRes = append(taxsRes, taxUpload)
	}

	return c.JSON(http.StatusOK, &TaxUploadResponse{Taxs: taxsRes})
}
