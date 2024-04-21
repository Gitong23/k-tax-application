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

func (p *Postgres) PersonalAllowance() (float64, error) {
	row, err := p.Db.Query("SELECT init_amount FROM allowances WHERE type = 'personal'")
	if err != nil {
		return 0, err
	}
	defer row.Close()

	var personalAllowance float64
	for row.Next() {
		err := row.Scan(&personalAllowance)
		if err != nil {
			return 0, err
		}
	}
	return personalAllowance, nil
}
