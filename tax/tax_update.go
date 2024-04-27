package tax

import "fmt"

func validateInitPersonalDeduction(s Storer, amount float64) error {

	p, err := s.PersonalAllowance()
	if err != nil {
		return fmt.Errorf("Internal Server Error")
	}

	if amount < p.MinAmount || amount > p.MaxAmount {
		return fmt.Errorf("Invalid personal deduction amount")
	}

	return nil
}
