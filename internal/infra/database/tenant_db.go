package database

import (
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"gorm.io/gorm"
)

type TenantDB struct {
	DB *gorm.DB
}

func NewTenantDB(db *gorm.DB) *TenantDB {
	return &TenantDB{DB: db}
}

func (tdb *TenantDB) Create(tenant *entity.Tenant) error {
	return tdb.DB.Create(tenant).Error
}

func (tdb *TenantDB) FindByID(id uint64) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := tdb.DB.Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (tdb *TenantDB) Update(tenant *entity.Tenant) error {
	_, err := tdb.FindByID(uint64(tenant.ID))
	if err != nil {
		return err
	}
	return tdb.DB.Save(tenant).Error
}

func (tdb *TenantDB) FindAll(page int, limit int, sort string) ([]entity.Tenant, error) {
	var tenants []entity.Tenant
	var err error
	if sort == "" || sort == "asc" || sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		err = tdb.DB.Limit(limit).Offset((page - 1) * limit).Order("id " + sort).Find(&tenants).Error
	} else {
		err = tdb.DB.Order("CreatedAt " + sort).Find(&tenants).Error
	}
	return tenants, err
}

func (tdb *TenantDB) Delete(id uint64) error {

    if err := tdb.DB.Where("tenant_id = ?", id).Delete(&entity.Project{}).Error; err != nil {
        return err
    }

    return tdb.DB.Where("id = ?", id).Delete(&entity.Tenant{}).Error
}
