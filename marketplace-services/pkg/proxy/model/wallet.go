package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/jinzhu/gorm"
)

type Wallet struct {
	gorm.Model
	AccountID  uint   `gorm:"unique" sql:"type:integer REFERENCES accounts(id)"`
	Passphrase string `gorm:"not null"`
	Address    []byte `gorm:"unique;not null"`
	FilePath   string `gorm:"not null"`
	PublicKey  []byte `gorm:"unique;not null"`
}

func (w Wallet) Validate() error {
	return validation.ValidateStruct(&w,
		validation.Field(&w.AccountID, validation.Required),
		validation.Field(&w.Passphrase, validation.Required, validation.Length(8, 32)),
		validation.Field(&w.Address, validation.Length(20, 20)),
	)
}
