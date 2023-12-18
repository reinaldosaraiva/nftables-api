// table_db.go
package database

import (
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"gorm.io/gorm"
)

type TableDB struct {
    DB *gorm.DB
}

func NewTableDB(db *gorm.DB) *TableDB {
    return &TableDB{DB: db}
}

func (tdb *TableDB) Create(table *entity.Table) error {
    return tdb.DB.Create(table).Error
}

func (tdb *TableDB) FindByID(id uint64) (*entity.Table, error) {
    var table entity.Table
    err := tdb.DB.Where("id = ?", id).First(&table).Error
    if err != nil {
        return nil, err
    }
    return &table, nil
}

func (tdb *TableDB) FindByName(name string) (*entity.Table, error) {
    var table entity.Table
    err := tdb.DB.Where("name = ?", name).First(&table).Error
    if err != nil {
        return nil, err
    }
    return &table, nil
}

func (tdb *TableDB) Update(table *entity.Table) error {
    
    _, err := tdb.FindByID(uint64(table.ID))
    if err != nil {
        return err
    }
    return tdb.DB.Save(table).Error
}

func (tdb *TableDB) FindAll(page int, limit int, sort string) ([]entity.Table, error) {
    var tables []entity.Table
    var err error
    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }
    if page != 0 && limit != 0 {
        err = tdb.DB.Limit(limit).Offset((page - 1) * limit).Order("id " + sort).Find(&tables).Error
    } else {
        err = tdb.DB.Order("CreatedAt " + sort).Find(&tables).Error
    }
    return tables, err
}

func (tdb *TableDB) Delete(id uint64) error {
    table, err := tdb.FindByID(id)
    if err != nil {
        return err
    }
    return tdb.DB.Delete(table).Error
}
