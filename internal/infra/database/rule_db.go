package database

import (
	"errors"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"gorm.io/gorm"
)

type RuleDB struct {
    DB *gorm.DB
}

func NewRuleDB(db *gorm.DB) *RuleDB {
    return &RuleDB{DB: db}
}

func (rdb *RuleDB) Create(rule *entity.Rule) error {
    tx := rdb.DB.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    // Verifique se a Chain associada existe
    if err := tx.First(&entity.Chain{}, rule.ChainID).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("chain not found")
        }
        return err
    }

    if err := tx.Create(rule).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

func (rdb *RuleDB) FindByID(id uint64) (*entity.Rule, error) {
    var rule entity.Rule
    err := rdb.DB.Preload("Chain").Where("id = ?", id).First(&rule).Error
    if err != nil {
        return nil, err
    }
    return &rule, nil
}

func (rdb *RuleDB) Update(rule *entity.Rule) error {
    _, err := rdb.FindByID(uint64(rule.ID))
    if err != nil {
        return err
    }
    return rdb.DB.Save(rule).Error
}

func (rdb *RuleDB) FindAll(page int, limit int, sort string) ([]entity.Rule, error) {
    var rules []entity.Rule
    var err error
    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }
    if page != 0 && limit != 0 {
        err = rdb.DB.Preload("Chain").Limit(limit).Offset((page - 1) * limit).Order("id " + sort).Find(&rules).Error
    } else {
        err = rdb.DB.Preload("Chain").Order("CreatedAt " + sort).Find(&rules).Error
    }
    return rules, err
}

func (rdb *RuleDB) Delete(id uint64) error {
    rule, err := rdb.FindByID(id)
    if err != nil {
        return err
    }
    return rdb.DB.Delete(rule).Error
}
