package entity

import "gorm.io/gorm"

type Project struct {
	Name        string `json:"name"`
	TenantID   uint64 `json:"tenant_id"`
	Tenant    	Tenant 
	// Chains      []Chain
	gorm.Model
}

func NewProject(name string, tenantID uint64) (*Project,error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if tenantID == 0 {
		return nil, ErrTenantRequired
	}
	return &Project{
		Name: name,
		TenantID: tenantID,
	}, nil

}
