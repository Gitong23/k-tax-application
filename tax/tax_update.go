package tax

import (
	"fmt"
)

func validateInitPersonalDeduction(s Storer, amount float64) (int, error) {

	p, err := s.PersonalAllowance()
	if err != nil {
		return 500, fmt.Errorf("Internal Server Error")
	}

	if amount < p.MinAmount || amount > p.MaxAmount {
		return 400, fmt.Errorf("Invalid personal deduction amount")
	}

	return 200, nil
}

func validateMaxKreceipt(s Storer, amount float64) (int, error) {

	k, err := s.KreceiptAllowance()
	if err != nil {
		return 500, fmt.Errorf("Internal Server Error")
	}

	if amount < k.MinAmount || amount > k.LimitMaxAmount {
		return 400, fmt.Errorf("Invalid K-receipt deduction amount")
	}

	return 200, nil
}
