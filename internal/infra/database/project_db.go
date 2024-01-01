package database

import (
	"errors"

	"github.com/reinaldosaraiva/nftables-api/internal/dto"
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
    var tenant entity.Tenant
    err := pdb.DB.First(&tenant, project.TenantID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("tenant not found")
        }
        return err
    }
    return pdb.DB.Create(project).Error
}

func (pdb *ProjectDB) FindByID(id uint64) (*dto.DetailsProjectDTO, error) {
    var projectDTO dto.DetailsProjectDTO
    err := pdb.DB.Table("projects").
        Select("projects.id, projects.name, tenants.id as tenant_id, tenants.name as tenant_name").
        Joins("left join tenants on tenants.id = projects.tenant_id").
        Where("projects.id = ?", id).Where("projects.deleted_at is null").
        Scan(&projectDTO).Error

    if err != nil {
        return nil, err
    }
    if projectDTO.ID == 0 {
        return nil, gorm.ErrRecordNotFound
    }
    return &projectDTO, nil
}




func (pdb *ProjectDB) FindByName(name string) (*dto.DetailsProjectDTO, error) {
    var projectDTO dto.DetailsProjectDTO
    err := pdb.DB.Table("projects").
        Select("projects.id, projects.name, tenants.id as tenant_id, tenants.name as tenant_name").
        Joins("left join tenants on tenants.id = projects.tenant_id").
        Where("projects.name = ?", name).Where("projects.deleted_at is null").
        Scan(&projectDTO).Error

    if err != nil {
        return nil, err
    }

    if projectDTO.ID == 0 {
        return nil, gorm.ErrRecordNotFound
    }

    return &projectDTO, nil
}

func (pdb *ProjectDB) Update(project *entity.Project) error {
    result := pdb.DB.Model(&entity.Project{}).Where("id = ?", project.ID).Updates(project)
    
    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }

    return nil
}


func (pdb *ProjectDB) FindAll(page int, limit int, sort string) ([]dto.DetailsProjectDTO, error) {
    var projectsDTO []dto.DetailsProjectDTO

    if sort == "" || sort == "asc" || sort != "desc" {
        sort = "asc"
    }

    query := pdb.DB.Table("projects").Select(
        "projects.*, tenants.id as tenant_id, tenants.name as tenant_name").
    Joins("left join tenants on tenants.id = projects.tenant_id").Where("projects.deleted_at is null")
    if page != 0 && limit != 0 {
        query = query.Limit(limit).Offset((page - 1) * limit)
    }
    if sort != "" {
        query = query.Order("projects.id " + sort)
    }
    
    err := query.Scan(&projectsDTO).Error
    if len(projectsDTO) == 0 {
        return []dto.DetailsProjectDTO{}, nil
    }

    return projectsDTO, err
}

func (pdb *ProjectDB) Delete(id uint64) error {
    tx := pdb.DB.Begin()
    if tx.Error != nil {
        return tx.Error
    }
    var project entity.Project
    if err := tx.First(&project, id).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return gorm.ErrRecordNotFound
        }
        return err
    }

    if err := tx.Delete(&project).Error; err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit().Error
}



