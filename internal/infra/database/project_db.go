package database

import (
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"gorm.io/gorm"
)

type ProjectDB struct {
    DB *gorm.DB
}

func NewProjectDB(db *gorm.DB) *ProjectDB {
    return &ProjectDB{DB: db}
}

func (pdb *ProjectDB) Create(project *entity.Project) error {
    return pdb.DB.Create(project).Error
}

func (pdb *ProjectDB) FindByID(id uint64) (*entity.Project, error) {
    var project entity.Project
    err := pdb.DB.Preload("Tenant").Preload("Chains").Where("id = ?", id).First(&project).Error
    if err != nil {
        return nil, err
    }
    return &project, nil
}

func (pdb *ProjectDB) FindByName(name string) (*entity.Project, error) {
    var project entity.Project
    err := pdb.DB.Preload("Tenant").Preload("Chains").Where("name = ?", name).First(&project).Error
    if err != nil {
        return nil, err
    }
    return &project, nil
}

func (pdb *ProjectDB) Update(project *entity.Project) error {
    _, err := pdb.FindByID(uint64(project.ID))
    if err != nil {
        return err
    }
    return pdb.DB.Save(project).Error
}

func (pdb *ProjectDB) FindAll(page int, limit int, sort string) ([]entity.Project, error) {
    var projects []entity.Project
    var err error
    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }
    if page != 0 && limit != 0 {
        err = pdb.DB.Preload("Tenant").Preload("Chains").Limit(limit).Offset((page - 1) * limit).Order("id " + sort).Find(&projects).Error
    } else {
        err = pdb.DB.Preload("Tenant").Preload("Chains").Order("CreatedAt " + sort).Find(&projects).Error
    }
    return projects, err
}

func (pdb *ProjectDB) Delete(id uint64) error {
    project, err := pdb.FindByID(id)
    if err != nil {
        return err
    }
    return pdb.DB.Delete(project).Error
}
