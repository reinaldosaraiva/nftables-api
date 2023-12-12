package entity

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNameRequired  = errors.New("Name is required")
	ErrDescriptionRequired = errors.New("Description is required")
	ErrTypeRequired = errors.New("Type is required")
	ErrStateRequired  = errors.New("State is required")
	ErrProjectRequired  = errors.New("Project is required")
	ErrTenantRequired = errors.New("Tenant is required")
)

type Table struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	State       string `json:"state"`
	Chains      []Chain `gorm:"foreignKey:TableID"`
	TenantID    uint64 `json:"tenant_id"`
	ProjectIDs []uint64   `json:"project_ids" gorm:"many2many:table_projects"`
	gorm.Model
}

func NewTable(name, description, type_name, state string) (*Table, error) {
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
	return &Table{
		Name:      name,
		Description: description,
		Type:      type_name,
		State:     state,
	}, nil
}
