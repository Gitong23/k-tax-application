package tax

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type Stub struct {
	err error
}

func (s *Stub) PersonalAllowance() (float64, error) {
	return 60000, s.err
}

func (s *Stub) DonationAllowance() (*Allowances, error) {
	return &Allowances{
		ID:             2,
		Type:           "donation",
		InitAmount:     0,
		MinAmount:      0,
		MaxAmount:      100000.0,
		LimitMaxAmount: 100000.0,
		CreatedAt:      "2024-04-22",
	}, s.err
}

func (s *Stub) KreceiptAllowance() (*Allowances, error) {
	return &Allowances{
		ID:             3,
		Type:           "k-receipt",
		InitAmount:     0,
		MinAmount:      0,
		MaxAmount:      50000.0,
		LimitMaxAmount: 100000.0,
		CreatedAt:      "2024-04-22",
	}, s.err
}

func TestCalTax(t *testing.T) {

	//TODO: Implement to test table
	tests := []struct {
		name     string
		reqBody  TaxRequest
		wantRes  Tax
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
			wantRes:  Tax{Tax: 0},
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
			wantRes:  Tax{Tax: 29000},
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
			wantRes:  Tax{Tax: 71000},
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
			wantRes:  Tax{Tax: 639000},
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
			wantRes:  Tax{Tax: 4000.0},
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
			wantRes:  Tax{Tax: 0.0},
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
			wantRes:  Tax{Tax: 0.0},
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
			wantRes:  Tax{Tax: 19000.0},
			wantHttp: http.StatusOK,
		},
		{
			name: "Income 500k wht 0 allowance donation 200k k-receipt 10k tax should be 18k",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        200000.0,
					},
					{
						AllowanceType: "k-receipt",
						Amount:        10000.0,
					},
				},
			},
			wantRes:  Tax{Tax: 18000.0},
			wantHttp: http.StatusOK,
		},
		{
			name: "Income 500k wht 0 allowance donation 200k k-receipt 100k tax should be 14k",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        200000.0,
					},
					{
						AllowanceType: "k-receipt",
						Amount:        100000.0,
					},
				},
			},
			wantRes:  Tax{Tax: 14000.0},
			wantHttp: http.StatusOK,
		},
		{
			name: "Income 500k wht 2k allowance donation 50k k-receipt 50k tax should be 17k",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         2000.0,
				Allowances: []AllowanceReq{
					{
						AllowanceType: "donation",
						Amount:        50000.0,
					},
					{
						AllowanceType: "k-receipt",
						Amount:        50000.0,
					},
				},
			},
			wantRes:  Tax{Tax: 17000.0},
			wantHttp: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			reqBodyStr, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Errorf("error marshalling json: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/tax/calculation", strings.NewReader(string(reqBodyStr)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/tax/calculation")

			stubTax := &Stub{err: nil}
			handler := NewHandler(stubTax)
			handler.CalTax(c)

			var got Tax

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
