package types

import "time"

type CreateOrderRequestModel struct {
	CustomerId      string      `json:"customer_id,omitempty" validate:"required,dive,required"`
	Items           []OrderItem `json:"items" validate:"required,dive,required"`
	Discounts       []*Discount `json:"discounts,omitempty"`
	ShippingAddress Address     `json:"shipping_address" validate:"required"`
	BillingAddress  Address     `json:"billing_address" validate:"required"`
}

type OrderResponseModel struct {
	Id              string      `json:"id"`
	CustomerId      string      `json:"customer_id"`
	Items           []OrderItem `json:"items"`
	ShippingAddress Address     `json:"shipping_address"`
	BillingAddress  Address     `json:"billing_address"`
	TotalPrice      float64     `json:"total_price"`
	Status          string      `json:"status"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Discounts       []*Discount `json:"discounts,omitempty"`
	IsDelete        bool        `json:"is_delete"`
}
type CustomerResponseModel struct {
	Id        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     []Phone   `json:"phone"`
	Address   []Address `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Pagination struct {
	Page  int
	Limit int
}

type OrderWithCustomerResponse struct {
	OrderResponseModel
	Customer CustomerResponseModel `json:"customer,omitempty"`
}

type OrderPriceInfo struct {
	TotalPrice float64   `bson:"total_price"`
	Discount   *Discount `bson:"discount,omitempty"`
}

type FinalPriceResult struct {
	OriginalPrice   float64 `json:"original_price"`
	DiscountApplied float64 `json:"discount_applied"`
	FinalPrice      float64 `json:"final_price"`
	DiscountType    string  `json:"discount_type,omitempty"`
}

type AggregationResult struct {
	TotalPrice float64    `bson:"total_price"`
	Discount   []Discount `bson:"discount"`
}
