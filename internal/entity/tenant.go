package entity

import (
	"gorm.io/gorm"
)


type Tenant struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Projects    []Project `gorm:"foreignKey:TenantID"`
	gorm.Model
}

func NewTenant(name string) (*Tenant, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	return &Tenant{
		Name: name,
	}, nil
}
