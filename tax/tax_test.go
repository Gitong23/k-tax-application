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

func TestCalTax(t *testing.T) {

	//TODO: Implement to test table
	tests := []struct {
		name    string
		reqBody TaxRequest
		want    Tax
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
			want: Tax{Tax: 0},
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
			want: Tax{Tax: 29000},
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
			want: Tax{Tax: 71000},
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
			want: Tax{Tax: 639000},
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
			want: Tax{Tax: 4000.0},
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

			handler := NewHandler()
			handler.CalTax(c)

			var got Tax
			if rec.Code != http.StatusOK {
				t.Errorf("expected status code %d but got %d", http.StatusOK, rec.Code)
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
