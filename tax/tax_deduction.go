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

func (a *Allowances) checkMin(amount float64) error {
	if amount < a.MinAmount {
		return fmt.Errorf("Invalid %s amount", a.Type)
	}
	return nil
}

type AllowanceChecker struct {
	a Allowances
}

func (d *Deductor) validateMin(each AllowanceReq) error {
	switch each.AllowanceType {
	case "donation":
		err := d.donation.checkMin(each.Amount)
		if err != nil {
			return err
		}

	case "k-receipt":
		err := d.kReceipt.checkMin(each.Amount)
		if err != nil {
			return err
		}
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
