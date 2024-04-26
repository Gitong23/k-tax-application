package tax

import (
	"fmt"
)

type Deductor struct {
	// personal Allowances
	// donation Allowances
	// kReceipt Allowances
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

func (d *Deductor) validateMin(each AllowanceReq) error {
	// m := d.mapAllowances[each.AllowanceType]
	if each.Amount < d.mapAllowances[each.AllowanceType].MinAmount {
		return fmt.Errorf("Invalid %s amount", each.AllowanceType)
	}
	return nil
}

func (d *Deductor) checkMinAllowanceReq(a []AllowanceReq) error {
	for _, each := range a {
		err := d.validateMin(each)
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
