package database

import (
	"errors"

	"github.com/reinaldosaraiva/nftables-api/internal/dto"
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
    var project entity.Project
    err := cdb.DB.First(&project, chain.ProjectID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("project not found")
        }
        return err
    }

    var table entity.Table
    err = cdb.DB.First(&table, chain.TableID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("table not found")
        }
        return err
    }
    return cdb.DB.Create(chain).Error
}

func (cdb *ChainDB) FindByID(id uint64) (*dto.DetailsChainDTO, error) {
    var chain entity.Chain
    err := cdb.DB.Preload("Project").Preload("Table").Preload("Rules").Where("id = ?", id).First(&chain).Error
    if err != nil {
        return nil, err
    }

    var rulesDTO []dto.DetailsRuleDTO
    for _, rule := range chain.Rules {
        ruleDTO := dto.DetailsRuleDTO{
            Protocol: rule.Protocol,
            Port:     rule.Port,
            Action:   rule.Action,
        }
        rulesDTO = append(rulesDTO, ruleDTO)
    }

    chainDTO := &dto.DetailsChainDTO{
        ID:         uint64(chain.ID),
        Name:        chain.Name,
        Type:        chain.Type,
        Priority:    chain.Priority,
        Policy:      chain.Policy,
        ProjectID:  uint64(chain.ProjectID),
        ProjectName: chain.Project.Name,
        TableID:    uint64(chain.TableID),
        TableName:   chain.Table.Name,
        Rules:       rulesDTO,
    }

    return chainDTO, nil
}

func (cdb *ChainDB) FindByName(name string) (*dto.DetailsChainDTO, error) {
    var chain entity.Chain
    err := cdb.DB.Preload("Project").Preload("Table").Preload("Rules").Where("name = ?", name).First(&chain).Error
    if err != nil {
        return nil, err
    }

    var rulesDTO []dto.DetailsRuleDTO
    for _, rule := range chain.Rules {
        ruleDTO := dto.DetailsRuleDTO{
            Protocol: rule.Protocol,
            Port:     rule.Port,
            Action:   rule.Action,
        }
        rulesDTO = append(rulesDTO, ruleDTO)
    }

    chainDTO := &dto.DetailsChainDTO{
        ID:         uint64(chain.ID),
        Name:        chain.Name,
        Type:        chain.Type,
        Priority:    chain.Priority,
        Policy:      chain.Policy,
        ProjectID:  uint64(chain.ProjectID),
        ProjectName: chain.Project.Name,
        TableID:    uint64(chain.TableID),
        TableName:   chain.Table.Name,
        Rules:       rulesDTO,
    }

    return chainDTO, nil
}

func (cdb *ChainDB) Update(chain *entity.Chain) error {
    _, err := cdb.FindByID(uint64(chain.ID))
    if err != nil {
        return err
    }
    return cdb.DB.Save(chain).Error
}

func (cdb *ChainDB) FindAll(page int, limit int, sort string) ([]dto.DetailsChainDTO, error) {
    var chains []entity.Chain
    var chainsDTO []dto.DetailsChainDTO

    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }

    query := cdb.DB.Model(&entity.Chain{})
    if page != 0 && limit != 0 {
        query = query.Limit(limit).Offset((page - 1) * limit)
    }
    if sort != "" {
        query = query.Order("id " + sort)
    }

    err := query.Preload("Project").Preload("Table").Preload("Rules").Find(&chains).Error
    if err != nil {
        return nil, err
    }

    for _, chain := range chains {
        var rulesDTO []dto.DetailsRuleDTO = make([]dto.DetailsRuleDTO, 0)
        for _, rule := range chain.Rules {
            ruleDTO := dto.DetailsRuleDTO{
                Protocol: rule.Protocol,
                Port:     rule.Port,
                Action:   rule.Action,
            }
            rulesDTO = append(rulesDTO, ruleDTO)
        }

        chainDTO := dto.DetailsChainDTO{
            ID:         uint64(chain.ID),
            Name:        chain.Name,
            Type:        chain.Type,
            Priority:    chain.Priority,
            Policy:      chain.Policy,
            ProjectID:  uint64(chain.ProjectID),
            ProjectName: chain.Project.Name,
            TableID:    uint64(chain.TableID),
            TableName:   chain.Table.Name,
            Rules:       rulesDTO,
        }

        chainsDTO = append(chainsDTO, chainDTO)
    }

    return chainsDTO, nil
}

// chain_db.go

func (cdb *ChainDB) Delete(id uint64) error {
    tx := cdb.DB.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    var chain entity.Chain
    if err := tx.First(&chain, id).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return gorm.ErrRecordNotFound
        }
        return err
    }

    if err := tx.Where("chain_id = ?", id).Delete(&entity.Rule{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Delete(&chain).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
