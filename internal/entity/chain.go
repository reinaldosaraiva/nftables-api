package entity

import (
	"errors"

	"gorm.io/gorm"
)

type Chain struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	State       string `json:"state"`
	ProjectID 	uint64  `json:"project_id"`
	Project	Project
	TableID 	uint64  `json:"table_id"`
	Table	Table
	Rules		[]Rule
	gorm.Model
}

func NewChain(name, description, type_name, state string, projectID uint64) (*Chain, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if description == "" {
		return nil, ErrDescriptionRequired
	}
	if type_name == "" {
		return nil, ErrTypeRequired
	}
	if state == "" {
		return nil, ErrStateRequired
	}
	if projectID == 0 {
		return nil, errors.New("Project is required")
	}
	return &Chain{
		Name:        name,
		Description: description,
		Type:        type_name,
		State:       state,
		ProjectID:   projectID,
	}, nil
}
