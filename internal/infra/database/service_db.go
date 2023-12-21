package database

import (
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"gorm.io/gorm"
)

type ServiceDB struct {
	DB *gorm.DB
}

func NewServiceDB(db *gorm.DB) *ServiceDB {
	return &ServiceDB{DB: db}
}

func (sdb *ServiceDB) Create(service *entity.Service) error {
	return sdb.DB.Create(service).Error
}

func (sdb *ServiceDB) FindByID(id uint64) (*entity.Service, error) {
	var service entity.Service
	err := sdb.DB.Preload("Rule").Where("id = ?", id).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (sdb *ServiceDB) FindByName(name string) (*entity.Service, error) {
	var service entity.Service
	err := sdb.DB.Preload("Rule").Where("name = ?", name).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (sdb *ServiceDB) Update(service *entity.Service) error {
	_, err := sdb.FindByID(uint64(service.ID))
	if err != nil {
		return err
	}
	return sdb.DB.Save(service).Error
}

func (sdb *ServiceDB) FindAll(page int, limit int, sort string) ([]entity.Service, error) {
	var services []entity.Service
	var err error
	if sort == "" || sort == "asc" || sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		err = sdb.DB.Preload("Rule").Limit(limit).Offset((page - 1) * limit).Order("id " + sort).Find(&services).Error
	} else {
		err = sdb.DB.Preload("Rule").Order("CreatedAt " + sort).Find(&services).Error
	}
	return services, err
}

func (sdb *ServiceDB) Delete(id uint64) error {
	service, err := sdb.FindByID(id)
	if err != nil {
		return err
	}
	return sdb.DB.Delete(service).Error
}