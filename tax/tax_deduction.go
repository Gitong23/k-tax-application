package tax

import (
	"fmt"
)

type Deductor struct {
	mapAllowances map[string]*Allowances
}

func NewDeductor(db Storer) (*Deductor, error) {

	personal, err := db.PersonalAllowance()
	if err != nil {
		return nil, err
	}

	donation, err := db.DonationAllowance()
	if err != nil {
		return nil, err
	}

	kReceipt, err := db.KreceiptAllowance()
	if err != nil {
		return nil, err
	}

	return &Deductor{
		mapAllowances: map[string]*Allowances{
			"personal":  personal,
			"donation":  donation,
			"k-receipt": kReceipt,
		},
	}, nil
}

func (d *Deductor) Min(t string) float64 {
	return d.mapAllowances[t].MinAmount
}

func (d *Deductor) validateMin(a float64, t string) error {
	if a < d.Min(t) {
		return fmt.Errorf("Invalid %s amount", t)
	}
	return nil
}

func (d *Deductor) checkMinAllowanceReq(a []AllowanceReq) error {
	for _, e := range a {
		err := d.validateMin(e.Amount, e.AllowanceType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Deductor) deductIncome(allReq []AllowanceReq) float64 {
	result := 0.0
	for _, allowance := range allReq {
		switch allowance.AllowanceType {
		case "donation":
			m := d.mapAllowances["donation"]
			if allowance.Amount > m.MaxAmount {
				result += m.MaxAmount
				break
			}

			result += allowance.Amount
		case "k-receipt":
			m := d.mapAllowances["k-receipt"]
			if allowance.Amount > m.MaxAmount {
				result += m.MaxAmount
				break
			}

			result += allowance.Amount
		default:
			result += 0
		}
	}

	m := d.mapAllowances["personal"]
	return result + m.InitAmount
}
