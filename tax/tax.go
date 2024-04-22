package tax

type AllowanceReq struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxRequest struct {
	TotalIncome float64        `json:"totalIncome"`
	WHT         float64        `json:"wht"`
	Allowances  []AllowanceReq `json:"allowances"`
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

type TaxResponse struct {
	Tax       float64     `json:"tax"`
	TaxRefund interface{} `json:"taxRefund,omitempty"`
}
