package internal

import (
	"tesodev-korpes/CustomerService/config"
	"tesodev-korpes/CustomerService/internal/types"

	"github.com/google/uuid"
)

func ToCustomerResponse(customer *types.Customer) *types.CustomerResponseModel {
	if customer == nil {
		return nil
	}
	return &types.CustomerResponseModel{
		ID:        customer.Id,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Address:   customer.Address,
		IsActive:  customer.IsActive,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
		Role:      customer.Role,
	}
}

func FromCreateCustomerRequest(req *types.CreateCustomerRequestModel) *types.Customer {
	if req == nil {
		return nil
	}

	addresses := make([]types.Address, len(req.Address))
	for i, addr := range req.Address {
		addresses[i] = types.Address{
			Id:      uuid.New().String(),
			City:    addr.City,
			State:   addr.State,
			ZipCode: addr.ZipCode,
		}
	}

	phones := make([]types.Phone, len(req.Phone))
	for i, ph := range req.Phone {
		phones[i] = types.Phone{
			Id:          uuid.New().String(),
			PhoneNumber: ph.PhoneNumber,
		}
	}

	return &types.Customer{
		Id:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     phones,
		Address:   addresses,
		Password:  req.Password,
		Role: types.Role{
			SystemRole: config.RoleStatus.System.NonPremium,
			Membership: config.RoleStatus.Membership.User,
		},
		IsActive: true,
	}

}
func FromUpdateCustomerRequest(customer *types.Customer, req *types.UpdateCustomerRequestModel) *types.Customer {
	if customer == nil || req == nil {
		return customer
	}
	if req.FirstName != "" {
		customer.FirstName = req.FirstName
	}
	if req.LastName != "" {
		customer.LastName = req.LastName
	}

	if req.Phone != nil {
		for i, p := range req.Phone {
			if p.Id == "" && i < len(customer.Phone) {
				req.Phone[i].Id = customer.Phone[i].Id
			}
		}
		customer.Phone = req.Phone
	}

	if req.Address != nil {
		for i, a := range req.Address {
			if a.Id == "" && i < len(customer.Address) {
				req.Address[i].Id = customer.Address[i].Id
			}
		}
		customer.Address = req.Address
	}

	customer.IsActive = req.IsActive
	return customer
}

func FromCustomerResponse(resp *types.CustomerResponseModel) *types.Customer {
	if resp == nil {
		return nil
	}
	return &types.Customer{
		Id:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		Phone:     resp.Phone,
		Address:   resp.Address,
		IsActive:  resp.IsActive,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
		Role:      resp.Role,
	}
}

func ToVerifiedUserFromResponse(c *types.CustomerResponseModel) types.VerifiedUser {
	return types.VerifiedUser{
		ID:    c.ID,
		Email: c.Email,
	}
}

func ToVerifyTokenResponse(user *types.CustomerResponseModel) types.VerifyTokenResponse {
	return types.VerifyTokenResponse{
		Message: "Token verified successfully",
		User:    ToVerifiedUserFromResponse(user),
	}
}

func ToLoginResponse(token string, customer *types.Customer) types.LoginResponse {
	return types.LoginResponse{
		Token:   token,
		User:    ToCustomerResponse(customer),
		Message: "Login successful",
	}
}
