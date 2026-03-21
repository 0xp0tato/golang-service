package util

import (
	"shared/enum"

	gofakeit "github.com/brianvoe/gofakeit/v7"
)

func CreateData() enum.Data {
	data := enum.Data{
		Name:        gofakeit.Name(),
		Email:       gofakeit.Email(),
		Phone:       gofakeit.Phone(),
		Company:     gofakeit.Company(),
		Ccn:         gofakeit.CreditCardNumber(nil),
		Designation: gofakeit.JobTitle(),
	}

	return data
}
