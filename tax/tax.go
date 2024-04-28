package tax

type AllowanceReq struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type Allowances struct {
	ID             int     `json:"id"`
	Type           string  `json:"type"`
	InitAmount     float64 `json:"init_amount"`
	MinAmount      float64 `json:"min_amount"`
	MaxAmount      float64 `json:"max_amount"`
	LimitMaxAmount float64 `json:"limit_max_amount"`
	CreatedAt      string  `json:"created_at"`
}

type TaxRequest struct {
	TotalIncome float64        `json:"totalIncome"`
	WHT         float64        `json:"wht"`
	Allowances  []AllowanceReq `json:"allowances"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxResponse struct {
	Tax       float64     `json:"tax"`
	TaxRefund interface{} `json:"taxRefund,omitempty"`
	TaxLevels []TaxLevel  `json:"taxLevels,omitempty"`
}

type DeductionReq struct {
	Amount float64 `json:"amount"`
}

type InitPersonalDeductRes struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type MaxKreceiptRes struct {
	Kreceipt float64 `json:"kReceipt"`
}

type TaxUpload struct {
	TotalIncome float64  `json:"totalIncome"`
	Tax         float64  `json:"tax"`
	TaxRefund   *float64 `json:"taxRefund,omitempty"`
}

type TaxUploadResponse struct {
	Taxs []TaxUpload `json:"taxs"`
}
