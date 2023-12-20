package entity

import (
	"gorm.io/gorm"
)

type NetworkObject struct {
	gorm.Model
	Name     string `json:"name"`
	Address  string `json:"address"`
}

func NewNetworkObject(name, address string) (*NetworkObject, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if address == "" {
		return nil, ErrAddressRequired
	}
	return &NetworkObject{
		Name: name,
		Address: address,
	}, nil
}
