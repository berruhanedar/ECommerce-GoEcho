package types

import (
	"tesodev-korpes/pkg/customError"
	"tesodev-korpes/pkg/validators"
)

func (c CreateOrderRequestModel) CreateValidate() *customError.AppError {

	if !validators.IsEmpty(c.ShippingAddress.City) {
		return customError.NewUnprocessableEntity(customError.InvalidAddressFormat, nil)
	}
	if !validators.IsEmpty(c.ShippingAddress.ZipCode) {
		return customError.NewUnprocessableEntity(customError.InvalidAddressFormat, nil)
	}
	if !validators.IsEmpty(c.ShippingAddress.State) {
		return customError.NewUnprocessableEntity(customError.InvalidAddressFormat, nil)
	}

	return nil
}
