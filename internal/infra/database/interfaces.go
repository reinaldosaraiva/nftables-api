package database

import "github.com/reinaldosaraiva/nftables-api/internal/entity"

type TableInterface interface {
    Create(table *entity.Table) error

    FindByID(id uint64) (*entity.Table, error)

    FindAll(page, limit int, sort string) ([]entity.Table, error)

    Update(table *entity.Table) error

    Delete(id uint64) error
}

type ChainInterface interface {
    Create(chain *entity.Chain) error

    FindByID(id uint64) (*entity.Chain, error)

    FindAll(page, limit int, sort string) ([]entity.Chain, error)

    Update(chain *entity.Chain) error

    Delete(id uint64) error
}

type RuleInterface interface {
    Create(rule *entity.Rule) error

    FindByID(id uint64) (*entity.Rule, error)

    FindAll(page, limit int, sort string) ([]entity.Rule, error)

    Update(rule *entity.Rule) error

    Delete(id uint64) error
}

type TenantInterface interface {
    Create(tenant *entity.Tenant) error

    FindByID(id uint64) (*entity.Tenant, error)

    FindAll(page, limit int, sort string) ([]entity.Tenant, error)

    Update(tenant *entity.Tenant) error

    Delete(id uint64) error
}

// ProjectInterface defines the interface for interacting with Project entities.
type ProjectInterface interface {
    // Create creates a new Project in the database.
    Create(project *entity.Project) error

    // FindByID finds a Project by its ID.
    FindByID(id uint64) (*entity.Project, error)

    // FindAll retrieves all Projects from the database.
    FindAll(page, limit int, sort string) ([]entity.Project, error)

    // Update updates a Project in the database.
    Update(project *entity.Project) error

    // Delete deletes a Project from the database.
    Delete(id uint64) error
}
