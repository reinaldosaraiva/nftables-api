package database

import (
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"gorm.io/gorm"
)

type NetworkObjectDB struct {
	DB *gorm.DB
}

func NewNetworkObjectDB(db *gorm.DB) *NetworkObjectDB {
	return &NetworkObjectDB{DB: db}
}

func (ndb *NetworkObjectDB) Create(networkObject *entity.NetworkObject) error {
	return ndb.DB.Create(networkObject).Error
}

func (ndb *NetworkObjectDB) FindByID(id uint64) (*entity.NetworkObject, error) {
	var networkObject entity.NetworkObject
	err := ndb.DB.Preload("Rule").Where("id = ?", id).First(&networkObject).Error
	if err != nil {
		return nil, err
	}
	return &networkObject, nil
}

func (ndb *NetworkObjectDB) FindByName(name string) (*entity.NetworkObject, error) {
	var networkObject entity.NetworkObject
	err := ndb.DB.Preload("Rule").Where("name = ?", name).First(&networkObject).Error
	if err != nil {
		return nil, err
	}
	return &networkObject, nil
}

func (ndb *NetworkObjectDB) Update(networkObject *entity.NetworkObject) error {
	_, err := ndb.FindByID(uint64(networkObject.ID))
	if err != nil {
		return err
	}
	return ndb.DB.Save(networkObject).Error
}

func (ndb *NetworkObjectDB) FindAll(page int, limit int, sort string) ([]entity.NetworkObject, error) {
	var networkObjects []entity.NetworkObject
	var err error
	if sort == "" || sort == "asc" || sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		err = ndb.DB.Preload("Rule").Limit(limit).Offset((page - 1) * limit).Order("id " + sort).Find(&networkObjects).Error
	} else {
		err = ndb.DB.Preload("Rule").Order("CreatedAt " + sort).Find(&networkObjects).Error
	}
	return networkObjects, err
}

func (ndb *NetworkObjectDB) Delete(id uint64) error {
	networkObject, err := ndb.FindByID(id)
	if err != nil {
		return err
	}
	return ndb.DB.Delete(networkObject).Error
}
