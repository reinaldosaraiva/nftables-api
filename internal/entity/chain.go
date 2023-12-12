package entity

import (
	"errors"

	"gorm.io/gorm"
)

type Chain struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	State       string `json:"state"`
	TableID    uint64 `json:"table_id"`
	Rules       []Rule `gorm:"foreignKey:ChainID"`
	ProjectID  uint64 `json:"project_id"`
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
