package tax

import (
	"fmt"
)

type Deductor struct {
	m map[string]*Allowances
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
		m: map[string]*Allowances{
			"personal":  personal,
			"donation":  donation,
			"k-receipt": kReceipt,
		},
	}, nil
}

func (d *Deductor) min(t string) float64 {
	return d.m[t].MinAmount
}

func (d *Deductor) max(t string) float64 {
	return d.m[t].MaxAmount
}

func (d *Deductor) initPer(t string) float64 {
	return d.m[t].InitAmount
}

func (d *Deductor) validateMin(a float64, t string) error {
	if a < d.min(t) {
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

func (d *Deductor) add(t string, a float64) float64 {
	if a > d.max(t) {
		return d.max(t)
	}
	return a
}

func (d *Deductor) deductIncome(a []AllowanceReq) float64 {
	result := 0.0
	for _, e := range a {
		result += d.add(e.AllowanceType, e.Amount)
	}

	return result + d.initPer("personal")
}

func (d *Deductor) checkMinMultiTaxReq(taxesReq []TaxRequest) error {
	for _, taxReq := range taxesReq {
		err := d.checkMinAllowanceReq(taxReq.Allowances)
		if err != nil {
			return err
		}
	}
	return nil
}
