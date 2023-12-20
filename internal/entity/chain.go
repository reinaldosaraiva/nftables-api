package entity

import (
	"gorm.io/gorm"
)

type Chain struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Priority	int    `json:"priority"`
	Policy		string `json:"policy"`
	ProjectID 	uint64  `json:"project_id"`
	Project	Project
	TableID 	uint64  `json:"table_id"`
	Table	Table
	Rules		[]Rule
	gorm.Model
}

func NewChain(name, type_name, policy string, priority int, projectID uint64, tableID uint64) (*Chain, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if policy == "" {
		return nil, ErrPolicyRequired
	}
	if type_name == "" {
		return nil, ErrTypeRequired
	}
	if tableID == 0 {
		return nil, ErrTableRequired
	}
	if projectID == 0 {
		return nil, ErrProjectRequired
	}
	return &Chain{
		Name:        name,
		Type:        type_name,
		Priority:    priority,
		Policy:       policy,
		ProjectID:   projectID,
		TableID:    tableID,
	}, nil
}
