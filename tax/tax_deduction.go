package tax

import (
	"fmt"
)

type Deductor struct {
	personal Allowances
	donation Allowances
	kReceipt Allowances
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
		personal: *personal,
		donation: *donation,
		kReceipt: *kReceipt,
	}, nil
}

func CheckMin(a Allowances, amount float64) error {
	if amount < a.MinAmount {
		return fmt.Errorf("Invalid %s amount", a.Type)
	}
	return nil
}

func (d *Deductor) Wtf(each AllowanceReq, amount float64) error {
	switch each.AllowanceType {
	case "donation":
		err := CheckMin(d.donation, each.Amount)
		if err != nil {
			return err
		}

	case "k-receipt":
		err := CheckMin(d.kReceipt, each.Amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Deductor) checkMinAllowanceReq(a []AllowanceReq) error {
	for _, each := range a {
		err := d.Wtf(each, each.Amount)
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
			if allowance.Amount > d.donation.MaxAmount {
				result += d.donation.MaxAmount
				break
			}

			result += allowance.Amount
		case "k-receipt":
			if allowance.Amount > d.kReceipt.MaxAmount {
				result += d.kReceipt.MaxAmount
				break
			}

			result += allowance.Amount
		default:
			result += 0
		}
	}
	return result + d.personal.InitAmount
}
