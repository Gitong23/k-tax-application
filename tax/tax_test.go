package tax

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Gitong23/assessment-tax/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Stub struct {
	personalAllowance *Allowances
	donationAllowance *Allowances
	kreceiptAllowance *Allowances
	adminUsername     string
	adminPassword     string
	err               error
}

func (s *Stub) PersonalAllowance() (*Allowances, error) {
	return s.personalAllowance, s.err
}

func (s *Stub) DonationAllowance() (*Allowances, error) {
	return s.donationAllowance, s.err
}

func (s *Stub) KreceiptAllowance() (*Allowances, error) {
	return s.kreceiptAllowance, s.err
}

func (s *Stub) UpdateInitPersonalAllowance(amount float64) error {
	s.personalAllowance.InitAmount = amount
	return s.err
}

func genTaxLevel(income float64) []TaxLevel {
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
			Tax:   s.taxStep(income),
		})
		income -= s.Max - s.Min
	}
	return taxLevels
}

func TestCalTax(t *testing.T) {

	//TODO: Implement to test table
	tests := []struct {
		name     string
		reqBody  TaxRequest
		wantRes  TaxResponse
		wantHttp int
	}{
		{
			name: "Income 120k wht 0 allowance 0 tax should be 0",
			reqBody: TaxRequest{
				TotalIncome: 120000.0,
				WHT:         0.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			wantRes:  TaxResponse{Tax: 0, TaxLevels: genTaxLevel(120000.0)},
			wantHttp: http.StatusOK,
		},
		{
			name: "Income 500k wht 0 allowance 0 tax should be 29000",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			wantRes:  TaxResponse{Tax: 29000, TaxLevels: genTaxLevel(440000.0)},
			wantHttp: http.StatusOK,
		},
		{
			name: "Income 800k wht 0 allowance 0 tax should be 71000",
			reqBody: TaxRequest{
				TotalIncome: 800000.0,
				WHT:         0.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			wantRes:  TaxResponse{Tax: 71000, TaxLevels: genTaxLevel(740000.0)},
			wantHttp: http.StatusOK,
		},
		{
			name: "Income 3M wht 0 allowance 0 tax should be 639000",
			reqBody: TaxRequest{
				TotalIncome: 3000000.0,
				WHT:         0.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			wantRes:  TaxResponse{Tax: 639000, TaxLevels: genTaxLevel(2940000)},
			wantHttp: http.StatusOK,
		},
		{
			name: "Income 500k wht 25k allowance 0 tax should be 4000",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         25000.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			wantRes:  TaxResponse{Tax: 4000.0, TaxLevels: genTaxLevel(440000.0)},
			wantHttp: http.StatusOK,
		},
		{
			name: "Wht can't be more than income",
			reqBody: TaxRequest{
				TotalIncome: 200000.0,
				WHT:         200001.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			wantRes:  TaxResponse{Tax: 0.0},
			wantHttp: http.StatusBadRequest,
		},
		{
			name: "Wht must more than 0",
			reqBody: TaxRequest{
				TotalIncome: 200000.0,
				WHT:         -5.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			wantRes:  TaxResponse{Tax: 0.0},
			wantHttp: http.StatusBadRequest,
		},
		{
			name: "Income 500k wht 0 allowance 200k tax should be 19000",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        200000.0,
					},
				},
			},
			wantRes:  TaxResponse{Tax: 19000.0, TaxLevels: genTaxLevel(340000.0)},
			wantHttp: http.StatusOK,
		},
		// {
		// 	name: "Income 500k wht 0 allowance donation 200k k-receipt 10k tax should be 18k",
		// 	reqBody: TaxRequest{
		// 		TotalIncome: 500000.0,
		// 		WHT:         0.0,
		// 		Allowances: []AllowanceReq{
		// 			{
		// 				AllowanceType: "donation",
		// 				Amount:        200000.0,
		// 			},
		// 			{
		// 				AllowanceType: "k-receipt",
		// 				Amount:        10000.0,
		// 			},
		// 		},
		// 	},
		// 	wantRes:  TaxResponse{Tax: 18000.0, TaxLevels: genTaxLevel(330000.0)},
		// 	wantHttp: http.StatusOK,
		// },
		// {
		// 	name: "Income 500k wht 0 allowance donation 200k k-receipt 100k tax should be 14k",
		// 	reqBody: TaxRequest{
		// 		TotalIncome: 500000.0,
		// 		WHT:         0.0,
		// 		Allowances: []AllowanceReq{
		// 			{
		// 				AllowanceType: "donation",
		// 				Amount:        200000.0,
		// 			},
		// 			{
		// 				AllowanceType: "k-receipt",
		// 				Amount:        100000.0,
		// 			},
		// 		},
		// 	},
		// 	wantRes:  TaxResponse{Tax: 14000.0, TaxLevels: genTaxLevel(290000.0)},
		// 	wantHttp: http.StatusOK,
		// },
		// {
		// 	name: "Income 500k wht 2k allowance donation 50k k-receipt 50k tax should be 17k",
		// 	reqBody: TaxRequest{
		// 		TotalIncome: 500000.0,
		// 		WHT:         2000.0,
		// 		Allowances: []AllowanceReq{
		// 			{
		// 				AllowanceType: "donation",
		// 				Amount:        50000.0,
		// 			},
		// 			{
		// 				AllowanceType: "k-receipt",
		// 				Amount:        50000.0,
		// 			},
		// 		},
		// 	},
		// 	wantRes:  TaxResponse{Tax: 17000.0, TaxLevels: genTaxLevel(340000.0)},
		// 	wantHttp: http.StatusOK,
		// },
		// {
		// 	name: "Minimum donation amount is 0",
		// 	reqBody: TaxRequest{
		// 		TotalIncome: 100000.0,
		// 		WHT:         0.0,
		// 		Allowances: []AllowanceReq{
		// 			{
		// 				AllowanceType: "donation",
		// 				Amount:        -1.0,
		// 			},
		// 		},
		// 	},
		// 	wantRes:  TaxResponse{Tax: 0.0},
		// 	wantHttp: http.StatusBadRequest,
		// },
		// {
		// 	name: "Minimum k-receipt amount is 0",
		// 	reqBody: TaxRequest{
		// 		TotalIncome: 100000.0,
		// 		WHT:         0.0,
		// 		Allowances: []AllowanceReq{
		// 			{
		// 				AllowanceType: "k-receipt",
		// 				Amount:        -50000.0,
		// 			},
		// 		},
		// 	},
		// 	wantRes:  TaxResponse{Tax: 0.0},
		// 	wantHttp: http.StatusBadRequest,
		// },
		// {
		// 	name: "Income 500k wht 20k allowance donation 200k k-receipt 10k tax should get refund 2k",
		// 	reqBody: TaxRequest{
		// 		TotalIncome: 500000.0,
		// 		WHT:         20000.0,
		// 		Allowances: []AllowanceReq{
		// 			{
		// 				AllowanceType: "donation",
		// 				Amount:        200000.0,
		// 			},
		// 			{
		// 				AllowanceType: "k-receipt",
		// 				Amount:        10000.0,
		// 			},
		// 		},
		// 	},
		// 	wantRes:  TaxResponse{Tax: 0.0, TaxLevels: genTaxLevel(330000.0), TaxRefund: 2000.0},
		// 	wantHttp: http.StatusOK,
		// },
	}

	stubTax := &Stub{
		personalAllowance: &Allowances{
			ID:             1,
			Type:           "personal",
			InitAmount:     60000,
			MinAmount:      0,
			MaxAmount:      100000.0,
			LimitMaxAmount: 100000.0,
			CreatedAt:      "2024-04-22",
		},
		donationAllowance: &Allowances{
			ID:             2,
			Type:           "donation",
			InitAmount:     0,
			MinAmount:      0,
			MaxAmount:      100000.0,
			LimitMaxAmount: 100000.0,
			CreatedAt:      "2024-04-22",
		},
		kreceiptAllowance: &Allowances{
			ID:             3,
			Type:           "k-receipt",
			InitAmount:     0,
			MinAmount:      0,
			MaxAmount:      50000.0,
			LimitMaxAmount: 100000.0,
			CreatedAt:      "2024-04-22",
		},
		err: nil,
	}

	e := echo.New()
	e.POST("/tax/calculations", NewHandler(stubTax).CalTax)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reqBodyStr, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Errorf("error marshalling json: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(string(reqBodyStr)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/tax/calculations")
			e.ServeHTTP(rec, req)

			var got TaxResponse

			if rec.Code != tt.wantHttp {
				t.Errorf("expected status code %d but got %d", tt.wantHttp, rec.Code)
			}

			err = json.Unmarshal(rec.Body.Bytes(), &got)
			if err != nil {
				t.Errorf("error unmarshalling json: %v", err)
			}

			if !reflect.DeepEqual(got, tt.wantRes) {
				t.Errorf("expected %v but got %v", tt.wantRes, got)
			}
		})
	}
}

func TestUpdatePersonalDeduction(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		httpWant int
		reqBody  DeductionReq
		wantRes  DeductionRes
	}{
		{
			name:     "Without Basic Auth",
			username: "user",
			password: "888",
			httpWant: http.StatusUnauthorized,
			reqBody: DeductionReq{
				Amount: 70000,
			},
			wantRes: DeductionRes{},
		},
		{
			name:     "Amount 70 k",
			username: "adminTax",
			password: "admin!",
			httpWant: http.StatusOK,
			reqBody: DeductionReq{
				Amount: 70000,
			},
			wantRes: DeductionRes{
				PersonalDeduction: 70000,
			},
		},
		{
			name:     "Exceed Max Amount",
			username: "adminTax",
			password: "admin!",
			httpWant: http.StatusBadRequest,
			reqBody: DeductionReq{
				Amount: 700000,
			},
			wantRes: DeductionRes{},
		},
		{
			name:     "Lower than Min Amount",
			username: "adminTax",
			password: "admin!",
			httpWant: http.StatusBadRequest,
			reqBody: DeductionReq{
				Amount: -50.0,
			},
			wantRes: DeductionRes{},
		},
	}

	stub := &Stub{
		personalAllowance: &Allowances{
			ID:             1,
			Type:           "personal",
			InitAmount:     60000,
			MinAmount:      0,
			MaxAmount:      100000.0,
			LimitMaxAmount: 100000.0,
			CreatedAt:      "2024-04-22",
		},
		adminUsername: "adminTax",
		adminPassword: "admin!",
		err:           nil,
	}

	e := echo.New()
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == stub.adminUsername && password == stub.adminPassword {
			return true, nil
		}
		return false, nil
	}))

	e.POST("/admin/deductions/personal", NewHandler(stub).SetPersonalDeduction)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reqBodyStr, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Errorf("error marshalling json: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", strings.NewReader(string(reqBodyStr)))
			req.SetBasicAuth(tt.username, tt.password)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/admin/deductions/personal")
			e.ServeHTTP(rec, req)

			if rec.Code != tt.httpWant {
				t.Errorf("expected status code %d but got %d", tt.httpWant, rec.Code)
			}

			var got DeductionRes
			err = json.Unmarshal(rec.Body.Bytes(), &got)
			if err != nil {
				t.Errorf("error unmarshalling json: %v", err)
			}

			if !reflect.DeepEqual(got, tt.wantRes) {
				t.Errorf("expected %v but got %v", tt.wantRes, got)
			}
		})
	}
}
