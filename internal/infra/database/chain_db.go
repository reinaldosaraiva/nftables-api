package database

import (
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"gorm.io/gorm"
)

type ChainDB struct {
    DB *gorm.DB
}

func NewChainDB(db *gorm.DB) *ChainDB {
    return &ChainDB{DB: db}
}

func (cdb *ChainDB) Create(chain *entity.Chain) error {
    return cdb.DB.Create(chain).Error
}

func (cdb *ChainDB) FindByID(id uint) (*entity.Chain, error) {
    var chain entity.Chain
    err := cdb.DB.Where("id = ?", id).First(&chain).Error
    if err != nil {
        return nil, err
    }
    return &chain, nil
}

func (cdb *ChainDB) Update(chain *entity.Chain) error {
    _, err := cdb.FindByID(chain.ID)
    if err != nil {
        return err
    }
    return cdb.DB.Save(chain).Error
}

func (cdb *ChainDB) FindAll(page int, limit int, sort string) ([]entity.Chain, error) {
    var chains []entity.Chain
    var err error
    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }
    if page != 0 && limit != 0 {
        err = cdb.DB.Limit(limit).Offset((page - 1) * limit).Order("id " + sort).Find(&chains).Error
    } else {
        err = cdb.DB.Order("CreatedAt " + sort).Find(&chains).Error
    }
    return chains, err
}

func (cdb *ChainDB) Delete(id uint) error {
    chain, err := cdb.FindByID(id)
    if err != nil {
        return err
    }
    return cdb.DB.Delete(chain).Error
}
