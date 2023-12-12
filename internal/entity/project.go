package entity

import "gorm.io/gorm"

type Project struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	TenantID    uint64 `json:"tenant_id"`
	Name        string `json:"name"`
	Tables      []Table `gorm:"many2many:table_projects"`
	Chains      []Chain `gorm:"foreignKey:ProjectID"`
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
