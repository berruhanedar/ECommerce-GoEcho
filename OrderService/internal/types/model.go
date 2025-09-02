package types

import "time"

type Discount struct {
	Id           string    `bson:"_id"`
	Role         string    `bson:"role"`
	StartDate    time.Time `bson:"start_date"`
	EndDate      time.Time `bson:"end_date"`
	DiscountCode string    `bson:"discount_code"`
	Type         string    `bson:"type"`
	Value        float64   `bson:"value"`
}

type Order struct {
	Id              string      `bson:"_id,omitempty"`
	CustomerId      string      `bson:"customer_id"`
	Items           []OrderItem `bson:"items"`
	ShippingAddress Address     `bson:"shipping_address"`
	BillingAddress  Address     `bson:"billing_address"`
	TotalPrice      float64     `bson:"total_price"`
	Discounts       []*Discount `bson:"discount,omitempty"`
	Status          string      `bson:"status"`
	CreatedAt       time.Time   `bson:"created_at"`
	UpdatedAt       time.Time   `bson:"updated_at"`
	IsDelete        bool        `bson:"is_delete"`
}

type OrderItem struct {
	ProductId   string  `bson:"product_id"`
	ProductName string  `bson:"product_name"`
	Quantity    int     `bson:"quantity"`
	UnitPrice   float64 `bson:"unit_price"`
}

type Address struct {
	Id      string `bson:"address_id,omitempty"`
	City    string `bson:"city"`
	State   string `bson:"state"`
	ZipCode string `bson:"zip_code"`
}

type Phone struct {
	Id          string `bson:"phone_id,omitempty"`
	PhoneNumber int    `bson:"phone_number"`
}


