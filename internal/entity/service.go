package entity

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name     string `json:"name"`
	Port    int    `json:"port"`
}

func NewService(name string, port int) (*Service, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if port == 0 {
		return nil, ErrPortRequired
	}
	return &Service{
			Name: name,
			Port: port,

	}, nil
}
