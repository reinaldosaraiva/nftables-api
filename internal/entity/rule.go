package entity

import (
	"errors"

	"gorm.io/gorm"
)

type Rule struct {
	ChainID    uint64 `json:"chain_id"`
	Chain	  Chain
	Port        int `json:"port"`
	Protocol    string `json:"protocol"`
	Action      string `json:"action"`
	ChainRules  []Chain `gorm:"many2many:rule_chain;"`
	ServiceRules []Service `gorm:"many2many:rule_service;"`
	NetworkObjectRules []NetworkObject `gorm:"many2many:rule_network_object;"`

	gorm.Model
}

func NewRule(chainID uint64,  protocol string, port int, action string) (*Rule, error) {
	if chainID == 0 {
		return nil, errors.New("Chain is required")
	}
	if port == 0 {
		return nil, errors.New("Port is required")
	}
	return &Rule{
		ChainID:    chainID,
		Protocol:   protocol,
		Port:       port,
		Action:     action,
	}, nil
}
