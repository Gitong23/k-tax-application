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
}

func TestCalTax(t *testing.T) {

	//TODO: Implement to test table
	tests := []struct {
		name       string
		reqBody    TaxRequest
		want       Tax
		httpStatus int
	}{
		{
			name: "Income 120k wht 0 allowance 0 tax should be 0",
			reqBody: TaxRequest{
				TotalIncome: 120000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			want:       Tax{Tax: 0},
			httpStatus: http.StatusOK,
		},
		{
			name: "Income 500k wht 0 allowance 0 tax should be 29000",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			want:       Tax{Tax: 29000},
			httpStatus: http.StatusOK,
		},
		{
			name: "Income 800k wht 0 allowance 0 tax should be 71000",
			reqBody: TaxRequest{
				TotalIncome: 800000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			want:       Tax{Tax: 71000},
			httpStatus: http.StatusOK,
		},
		{
			name: "Income 3M wht 0 allowance 0 tax should be 639000",
			reqBody: TaxRequest{
				TotalIncome: 3000000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			want:       Tax{Tax: 639000},
			httpStatus: http.StatusOK,
		},
		{
			name: "Income 500k wht 25k allowance 0 tax should be 4000",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         25000.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			want:       Tax{Tax: 4000.0},
			httpStatus: http.StatusOK,
		},
		{
			name: "Wht can't be more than income",
			reqBody: TaxRequest{
				TotalIncome: 200000.0,
				WHT:         200001.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			want:       Tax{Tax: 0.0},
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "Wht must more than 0",
			reqBody: TaxRequest{
				TotalIncome: 200000.0,
				WHT:         -5.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			want:       Tax{Tax: 0.0},
			httpStatus: http.StatusBadRequest,
		},
		{
			name: "Income 500k wht 0 allowance 200k tax should be 19000",
			reqBody: TaxRequest{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        200000.0,
					},
				},
			},
			want:       Tax{Tax: 19000.0},
			httpStatus: http.StatusOK,
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

			stubTax := &Stub{}
			handler := NewHandler(stubTax)
			handler.CalTax(c)

			var got Tax

			if rec.Code != tt.httpStatus {
				t.Errorf("expected status code %d but got %d", tt.httpStatus, rec.Code)
			}

			err = json.Unmarshal(rec.Body.Bytes(), &got)
			if err != nil {
				t.Errorf("error unmarshalling json: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}
