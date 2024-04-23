package postgres

import "github.com/Gitong23/assessment-tax/tax"

type Allowances struct {
	ID             int     `posgres:"id"`
	Type           string  `posgres:"type"`
	InitAmount     float64 `posgres:"init_amount"`
	MinAmount      float64 `posgres:"min_amount"`
	MaxAmount      float64 `posgres:"max_amount"`
	LimitMaxAmount float64 `posgres:"limit_max_amount"`
	CreatedAt      string  `posgres:"created_at"`
}

func (p *Postgres) PersonalAllowance() (*tax.Allowances, error) {
	row, err := p.Db.Query("SELECT * FROM allowances WHERE type = 'personal'")
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var personal tax.Allowances
	for row.Next() {
		err := row.Scan(
			&personal.ID,
			&personal.Type,
			&personal.InitAmount,
			&personal.MinAmount,
			&personal.MaxAmount,
			&personal.LimitMaxAmount,
			&personal.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return &personal, nil
}

func (p *Postgres) DonationAllowance() (*tax.Allowances, error) {
	row, err := p.Db.Query("SELECT * FROM allowances WHERE type = 'donation'")

	if err != nil {
		return nil, err
	}
	defer row.Close()

	var d tax.Allowances
	for row.Next() {
		err := row.Scan(
			&d.ID,
			&d.Type,
			&d.InitAmount,
			&d.MinAmount,
			&d.MaxAmount,
			&d.LimitMaxAmount,
			&d.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return &d, nil
}

func (p *Postgres) KreceiptAllowance() (*tax.Allowances, error) {
	row, err := p.Db.Query("SELECT * FROM allowances WHERE type = 'k-receipt'")

	if err != nil {
		return nil, err
	}
	defer row.Close()

	var k tax.Allowances
	for row.Next() {
		err := row.Scan(
			&k.ID,
			&k.Type,
			&k.InitAmount,
			&k.MinAmount,
			&k.MaxAmount,
			&k.LimitMaxAmount,
			&k.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return &k, nil
}

func (p *Postgres) UpdateInitPersonalAllowance(amount float64) error {
	_, err := p.Db.Exec("UPDATE allowances SET init_amount = $1 WHERE type = personal", amount)
	if err != nil {
		return err
	}

	return nil
}
