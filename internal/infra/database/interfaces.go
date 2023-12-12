package database

import "github.com/reinaldosaraiva/nftables-api/internal/entity"

// Interfaces for interacting with entity models.

// TableInterface defines the interface for interacting with Table entities.
type TableInterface interface {
    // Create creates a new Table in the database.
    Create(table *entity.Table) error

    // FindByID finds a Table by its ID.
    FindByID(id uint64) (*entity.Table, error)

    // FindAll retrieves all Tables from the database.
    FindAll(page, limit int, sort string) ([]entity.Table, error)

    // Update updates a Table in the database.
    Update(table *entity.Table) error

    // Delete deletes a Table from the database.
    Delete(id uint64) error
}

// ChainInterface defines the interface for interacting with Chain entities.
type ChainInterface interface {
    // Create creates a new Chain in the database.
    Create(chain *entity.Chain) error

    // FindByID finds a Chain by its ID.
    FindByID(id uint64) (*entity.Chain, error)

    // FindAll retrieves all Chains from the database.
    FindAll(page, limit int, sort string) ([]entity.Chain, error)

    // Update updates a Chain in the database.
    Update(chain *entity.Chain) error

    // Delete deletes a Chain from the database.
    Delete(id uint64) error
}

// RuleInterface defines the interface for interacting with Rule entities.
type RuleInterface interface {
    // Create creates a new Rule in the database.
    Create(rule *entity.Rule) error

    // FindByID finds a Rule by its ID.
    FindByID(id uint64) (*entity.Rule, error)

    // FindAll retrieves all Rules from the database.
    FindAll(page, limit int, sort string) ([]entity.Rule, error)

    // Update updates a Rule in the database.
    Update(rule *entity.Rule) error

    // Delete deletes a Rule from the database.
    Delete(id uint64) error
}

// TenantInterface defines the interface for interacting with Tenant entities.
type TenantInterface interface {
    // Create creates a new Tenant in the database.
    Create(tenant *entity.Tenant) error

    // FindByID finds a Tenant by its ID.
    FindByID(id uint64) (*entity.Tenant, error)

    // FindAll retrieves all Tenants from the database.
    FindAll(page, limit int, sort string) ([]entity.Tenant, error)

    // Update updates a Tenant in the database.
    Update(tenant *entity.Tenant) error

    // Delete deletes a Tenant from the database.
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
