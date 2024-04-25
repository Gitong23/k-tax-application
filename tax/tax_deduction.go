package tax

import "fmt"

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

func (d *Deductor) checkMinAllowanceReq(allReq []AllowanceReq) error {
	for _, allowance := range allReq {
		switch allowance.AllowanceType {
		case "donation":
			if allowance.Amount < d.donation.MinAmount {
				return fmt.Errorf("Invalid donation amount")
			}

		case "k-receipt":
			if allowance.Amount < d.kReceipt.MinAmount {
				return fmt.Errorf("Invalid k-receipt amount")
			}
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

