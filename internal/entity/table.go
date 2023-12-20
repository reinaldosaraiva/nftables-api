package entity

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNameRequired  = errors.New("Name is required")
	ErrDescriptionRequired = errors.New("Description is required")
	ErrTypeRequired = errors.New("Type is required")
	ErrStateRequired = errors.New("State is required")
	ErrCommentRequired  = errors.New("Comment is required")
	ErrProjectRequired  = errors.New("Project is required")
	ErrTenantRequired = errors.New("Tenant is required")
	ErrPolicyRequired = errors.New("Policy is required")
	ErrPriorityRequired = errors.New("Priority is required")
	ErrTableRequired = errors.New("Table is required")
	ErrPortRequired = errors.New("Port is required")
	ErrAddressRequired = errors.New("Address is required")
)

type Table struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Chains      []Chain 
	gorm.Model
}

func NewTable(name, description, type_name string,) (*Table, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if description == "" {
		return nil, ErrDescriptionRequired
	}
	if type_name == "" {
		return nil, ErrTypeRequired
	}

	return &Table{
		Name:      name,
		Description: description,
		Type:      type_name,		
	}, nil
}
