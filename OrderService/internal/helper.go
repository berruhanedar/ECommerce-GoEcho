package internal

import (
	"tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal/types"
	"time"

	"github.com/google/uuid"
)

func FromCreateOrderRequest(req *types.CreateOrderRequestModel) *types.Order {
	if req == nil {
		return nil
	}

	items := make([]types.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = types.OrderItem{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}
	discounts := make([]*types.Discount, len(req.Discounts))
	for i, d := range req.Discounts {
		if d != nil {
			discounts[i] = &types.Discount{
				Id:           d.Id,
				Role:         d.Role,
				DiscountCode: d.DiscountCode,
				Type:         d.Type,
				Value:        d.Value,
				StartDate:    d.StartDate,
				EndDate:      d.EndDate}
		}
	}

	return &types.Order{
		Id:              uuid.NewString(),
		CustomerId:      req.CustomerId,
		Items:           items,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		TotalPrice:      calculateTotalPrice(items),
		Status:          config.OrderStatus.Ordered,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Discounts:       discounts,
	}
}

func ToOrderResponse(order *types.Order) *types.OrderResponseModel {
	if order == nil {
		return nil
	}

	responseDiscounts := make([]*types.Discount, len(order.Discounts))
	for i, d := range order.Discounts {
		if d != nil {
			responseDiscounts[i] = &types.Discount{
				Id:           uuid.NewString(),
				Role:         d.Role,
				DiscountCode: d.DiscountCode,
				Type:         d.Type,
				Value:        d.Value,
				StartDate:    d.StartDate,
				EndDate:      d.EndDate,
			}
		}
	}

	return &types.OrderResponseModel{
		Id:              order.Id,
		CustomerId:      order.CustomerId,
		Items:           order.Items,
		ShippingAddress: order.ShippingAddress,
		BillingAddress:  order.BillingAddress,
		TotalPrice:      order.TotalPrice,
		Status:          order.Status,
		Discounts:       responseDiscounts,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func ToOrderWithCustomerResponse(order *types.OrderResponseModel, customer *types.CustomerResponseModel) *types.OrderWithCustomerResponse {
	if order == nil || customer == nil {
		return nil
	}

	return &types.OrderWithCustomerResponse{
		OrderResponseModel: *order,
		Customer:           *customer,
	}
}
