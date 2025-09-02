package types

import (
	"tesodev-korpes/pkg/customError"
	"tesodev-korpes/pkg/validators"
)

func (c CreateCustomerRequestModel) CreateValidate() *customError.AppError {

	if !validators.IsValidEmail(c.Email) {
		return customError.NewUnprocessableEntity(customError.InvalidEmailFormat, nil)
	}

	if !validators.IsValidName(c.FirstName) {
		return customError.NewUnprocessableEntity(customError.InvalidFirstName, nil)
	}

	if !validators.IsValidName(c.LastName) {
		return customError.NewUnprocessableEntity(customError.InvalidLastName, nil)
	}

	if !validators.IsValidPassword(c.Password) {
		return customError.NewUnprocessableEntity(customError.InvalidPasswordFormat, nil)
	}

	for _, addr := range c.Address {
		if !validators.IsEmpty(addr.City) ||
			!validators.IsEmpty(addr.ZipCode) ||
			!validators.IsEmpty(addr.State) {
			return customError.NewUnprocessableEntity(customError.InvalidAddressFormat, nil)
		}
	}

	for _, phone := range c.Phone {
		if !validators.IsValidPhone(phone.PhoneNumber) {
			return customError.NewUnprocessableEntity(customError.InvalidPhoneFormat, nil)
		}
	}

	return nil
}

func (c LoginRequestModel) LoginValidate() *customError.AppError {

	if !validators.IsValidEmail(c.Email) {
		return customError.NewUnprocessableEntity(customError.InvalidEmailFormat, nil)
	}

	if !validators.IsValidPassword(c.Password) {
		return customError.NewUnprocessableEntity(customError.InvalidPasswordFormat, nil)
	}

	return nil
}
