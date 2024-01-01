package database

import (
	"fmt"

	"github.com/reinaldosaraiva/nftables-api/internal/dto"
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

func (tdb *TenantDB) FindByID(id uint64) (*dto.CreateTenantDTO, error) {
	var tenantDTO dto.CreateTenantDTO
	err := tdb.DB.Model(&entity.Tenant{}).
		Where("tenants.id = ?", id).Where("tenants.deleted_at is null").
		First(&tenantDTO).
		Error
	
	if tenantDTO.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return &tenantDTO, nil
}

func (tdb *TenantDB) FindByName(name string) (*dto.CreateTenantDTO, error) {
	var tenantDTO dto.CreateTenantDTO
	err := tdb.DB.Model(&entity.Tenant{}).
		Where("tenants.name = ?", name).
		First(&tenantDTO).Where("tenants.deleted_at is null").
		Error
	if err != nil {
		return nil, err
	}
	return &tenantDTO, nil
}

func (tdb *TenantDB) Update(tenant *entity.Tenant) error {
	result, err := tdb.FindByID(uint64(tenant.ID))
	if err != nil {
		return err
	}
	fmt.Println(result)
	if result == nil {
		return gorm.ErrRecordNotFound
	}
	return tdb.DB.Save(tenant).Error
}

func (tdb *TenantDB) FindAll(page int, limit int, sort string) ([]dto.CreateTenantDTO, error) {
    var tenantsDTO []dto.CreateTenantDTO
    var err error

    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }

    query := tdb.DB.Model(&entity.Tenant{})
    if page != 0 && limit != 0 {
        query = query.Limit(limit).Offset((page - 1) * limit)
    }
    if sort != "" {
        query = query.Order("id " + sort)
    }
    
    err = query.Find(&tenantsDTO).Error
	if len(tenantsDTO) == 0 {
		return []dto.CreateTenantDTO{}, err
	}
    return tenantsDTO, err
}


func (tdb *TenantDB) Delete(id uint64) error {
    tx := tdb.DB.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    var projects []entity.Project
    if err := tx.Where("tenant_id = ?", id).Find(&projects).Error; err != nil {
        tx.Rollback()
        return err
    }

    for _, project := range projects {
        if err := tx.Where("project_id = ?", project.ID).Delete(&entity.Chain{}).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    if err := tx.Where("tenant_id = ?", id).Delete(&entity.Project{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Where("id = ?", id).Delete(&entity.Tenant{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

