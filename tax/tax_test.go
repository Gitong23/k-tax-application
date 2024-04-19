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
	t.Run("Income 120k allowance 0 tax should be 0", func(t *testing.T) {
		e := echo.New()
		reqBody := TaxRequest{
			TotalIncome: 120000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			}}
		reqBodyStr, err := json.Marshal(reqBody)
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

		want := Tax{Tax: 0}
		gotJson := rec.Body.Bytes()

		var got Tax

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, rec.Code)
		}

		err = json.Unmarshal(gotJson, &got)
		if err != nil {
			t.Errorf("error unmarshalling json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}

	})

	t.Run("Income 500k allowance 0 tax should be 29k", func(t *testing.T) {
		e := echo.New()
		reqBody := TaxRequest{
			TotalIncome: 500000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			}}
		reqBodyStr, err := json.Marshal(reqBody)
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

		want := Tax{Tax: 29000}
		gotJson := rec.Body.Bytes()

		var got Tax

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, rec.Code)
		}

		err = json.Unmarshal(gotJson, &got)
		if err != nil {
			t.Errorf("error unmarshalling json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}

	})

}
