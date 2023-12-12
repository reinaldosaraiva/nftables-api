package entity

import (
	"gorm.io/gorm"
)


type Tenant struct {
	Name        string `json:"name"`
	Projects    []Project 
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
