// table_db.go
package database

import (
	"errors"

	"github.com/reinaldosaraiva/nftables-api/internal/dto"
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

func (tdb *TableDB) FindByID(id uint64) (*dto.DetailsTableDTO, error) {
    var table entity.Table
    err := tdb.DB.Preload("Chains").Where("id = ?", id).First(&table).Error
    if err != nil {
        return nil, err
    }

    // Transforma as Chains da entidade em Chains do DTO
    chainsDTO := make([]dto.DetailsChainDTO, 0)
    for _, chain := range table.Chains {
        chainDTO := dto.DetailsChainDTO{
            Name: chain.Name,
            Type: chain.Type,
            Priority: chain.Priority,
            Policy: chain.Policy,
        }
        chainsDTO = append(chainsDTO, chainDTO)
    }

    tableDTO := &dto.DetailsTableDTO{
        ID:         uint64(table.ID),
        Name:        table.Name,
        Description: table.Description,
        Type:        table.Type,
        Chains:      chainsDTO,
    }

    return tableDTO, nil
}

func (tdb *TableDB) FindByName(name string) (*dto.DetailsTableDTO, error) {
    var table entity.Table
    err := tdb.DB.Preload("Chains").Where("name = ?", name).First(&table).Error
    if err != nil {
        return nil, err
    }

    chainsDTO := make([]dto.DetailsChainDTO, 0)
    for _, chain := range table.Chains {
        chainDTO := dto.DetailsChainDTO{
            Name: chain.Name,
            Type: chain.Type,
            Priority: chain.Priority,
            Policy: chain.Policy,
        }
        chainsDTO = append(chainsDTO, chainDTO)
    }

    tableDTO := &dto.DetailsTableDTO{
        ID:         uint64(table.ID),
        Name:        table.Name,
        Description: table.Description,
        Type:        table.Type,
        Chains:      chainsDTO,
    }

    return tableDTO, nil
}


func (tdb *TableDB) Update(table *entity.Table) error {
    
    _, err := tdb.FindByID(uint64(table.ID))
    if err != nil {
        return err
    }
    return tdb.DB.Save(table).Error
}

func (tdb *TableDB) FindAll(page int, limit int, sort string) ([]dto.DetailsTableDTO, error) {
    var tables []entity.Table
    var tablesDTO []dto.DetailsTableDTO

    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }

    query := tdb.DB.Model(&entity.Table{})
    if page != 0 && limit != 0 {
        query = query.Limit(limit).Offset((page - 1) * limit)
    }
    if sort != "" {
        query = query.Order("id " + sort)
    }

    err := query.Find(&tables).Error
    if err != nil {
        return nil, err
    }

    for _, table := range tables {
        var chains []entity.Chain
        tdb.DB.Where("table_id = ?", table.ID).Find(&chains)

        var chainsDTO []dto.DetailsChainDTO = make([]dto.DetailsChainDTO, 0)
        for _, chain := range chains {
            chainsDTO = append(chainsDTO, dto.DetailsChainDTO{
                Name: chain.Name,
                Type: chain.Type,
                Priority: chain.Priority,
                Policy: chain.Policy,
            })
        }

        tableDTO := dto.DetailsTableDTO{
            ID:         uint64(table.ID),
            Name:        table.Name,
            Description: table.Description,
            Type:        table.Type,
            Chains:      chainsDTO,
        }

        tablesDTO = append(tablesDTO, tableDTO)
    }

    return tablesDTO, nil
}

func (tdb *TableDB) Delete(id uint64) error {
    tx := tdb.DB.Begin()
    if tx.Error != nil {
        return tx.Error
    }
    var table entity.Table
    if err := tx.First(&table, id).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return gorm.ErrRecordNotFound
        }
        return err
    }
    if err := tx.Where("table_id = ?", id).Delete(&entity.Chain{}).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Delete(&table).Error; err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit().Error
}

