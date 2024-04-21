package postgres

type Allowances struct {
	ID             int     `posgres:"id"`
	Type           string  `posgres:"type"`
	InitAmount     float64 `posgres:"init_amount"`
	MinAmount      float64 `posgres:"min_amount"`
	MaxAmount      float64 `posgres:"max_amount"`
	LimitMaxAmount float64 `posgres:"limit_max_amount"`
	CreatedAt      string  `posgres:"created_at"`
}
