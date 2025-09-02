package types

import (
	"time"
)

type Customer struct {
	Id        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Password  string    `bson:"password" json:"password"`
	FirstName string    `bson:"first_name" json:"first_name"`
	LastName  string    `bson:"last_name" json:"last_name"`
	Email     string    `bson:"email" json:"email"`
	Role      Role      `bson:"role" json:"role"`
	Phone     []Phone   `bson:"phone" json:"phone"`
	Address   []Address `bson:"address" json:"address"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	IsActive  bool      `bson:"is_active" json:"is_active"`
	Token     string    `bson:"token" json:"token"`
}

type Address struct {
	Id      string `bson:"address_id,omitempty" json:"address_id,omitempty"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	ZipCode string `bson:"zip_code" json:"zip_code"`
}

type Phone struct {
	Id          string `bson:"phone_id,omitempty" json:"phone_id,omitempty"`
	PhoneNumber int    `bson:"phone_number" json:"phone_number"`
}

type Role struct {
	SystemRole string `bson:"role"`
	Membership string `bson:"membership"`
}
