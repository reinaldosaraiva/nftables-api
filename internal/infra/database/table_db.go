package database

import (
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"gorm.io/gorm"
)

type Table struct {
    DB *gorm.DB
}

func NewTable(db *gorm.DB) *Table {
    return &Table{DB: db}
}

func (t *Table) Create(table *entity.Table) error {
    return t.DB.Create(table).Error
}

func (t *Table) FindByID(id uint) (*entity.Table, error) {
    var table entity.Table
    err := t.DB.Where("id = ?", id).First(&table).Error
    if err != nil {
        return nil, err
    }
    return &table, nil
}

func (t *Table) Update(table *entity.Table) error {
    _, err := t.FindByID(table.ID)
    if err != nil {
        return err
    }
    return t.DB.Save(table).Error
}

func (t *Table) FindAll(page int, limit int, sort string) ([]entity.Table, error) {
    var tables []entity.Table
    var err error
    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }
    if page != 0 && limit != 0 {
        err = t.DB.Limit(limit).Offset((page - 1) * limit).Order("id " + sort).Find(&tables).Error
    } else {
        err = t.DB.Order("CreatedAt " + sort).Find(&tables).Error
    }
    return tables, err
}

func (t *Table) Delete(id uint) error {
    table, err := t.FindByID(id)
    if err != nil {
        return err
    }
    return t.DB.Delete(table).Error
}
