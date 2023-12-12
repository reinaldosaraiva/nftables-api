package entity

import (
	"errors"

	"gorm.io/gorm"
)

type Rule struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	ChainID    uint64 `json:"chain_id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Port        int `json:"port"`
	Protocol    string `json:"protocol"`
	Action      string `json:"action"`
	gorm.Model
}

func NewRule(chainID uint64, source, destination, protocol string, port int, action string) (*Rule, error) {
	if chainID == 0 {
		return nil, errors.New("Chain is required")
	}
	if port == 0 {
		return nil, errors.New("Port is required")
	}
	return &Rule{
		ChainID:    chainID,
		Source:     source,
		Destination: destination,
		Protocol:   protocol,
		Port:       port,
		Action:     action,
	}, nil
}
