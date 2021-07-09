package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Password []byte `gorm:"not null"`
	Role     Role   `gorm:"not null"`
	Wallet   Wallet `gorm:"association_autoupdate:false;association_autocreate:false"`
}

func (u Account) HasRole(role Role) bool {
	return u.Role == role
}

func (u Account) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(2, 32)),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 32)),
		validation.Field(&u.Role, validation.Required, validation.Min(1), validation.Max(2)),
	)
}
