package database

import "github.com/reinaldosaraiva/nftables-api/internal/entity"

type TableInterface interface {
    Create(table *entity.Table) error

    FindByID(id uint64) (*entity.Table, error)

    FindByName(name string) (*entity.Table, error)

    FindAll(page, limit int, sort string) ([]entity.Table, error)

    Update(table *entity.Table) error

    Delete(id uint64) error
}

type ChainInterface interface {
    Create(chain *entity.Chain) error

    FindByID(id uint64) (*entity.Chain, error)

    FindByName(name string) (*entity.Chain, error)

    FindAll(page, limit int, sort string) ([]entity.Chain, error)

    Update(chain *entity.Chain) error

    Delete(id uint64) error
}

type RuleInterface interface {
    Create(rule *entity.Rule) error

    FindByID(id uint64) (*entity.Rule, error)

    FindByName(name string) (*entity.Rule, error)

    FindAll(page, limit int, sort string) ([]entity.Rule, error)

    Update(rule *entity.Rule) error

    Delete(id uint64) error
}

type TenantInterface interface {
    Create(tenant *entity.Tenant) error

    FindByID(id uint64) (*entity.Tenant, error)

    FindByName(name string) (*entity.Tenant, error)

    FindAll(page, limit int, sort string) ([]entity.Tenant, error)

    Update(tenant *entity.Tenant) error

    Delete(id uint64) error
}

type ProjectInterface interface {
    Create(project *entity.Project) error

    FindByID(id uint64) (*entity.Project, error)

    FindByName(name string) (*entity.Project, error)

    FindAll(page, limit int, sort string) ([]entity.Project, error)

    Update(project *entity.Project) error

    Delete(id uint64) error
}

type ServiceInterface interface {
    Create(service *entity.Service) error

    FindByID(id uint64) (*entity.Service, error)

    FindByName(name string) (*entity.Service, error)

    FindAll(page, limit int, sort string) ([]entity.Service, error)

    Update(service *entity.Service) error

    Delete(id uint64) error
}

type NetworkObjectInterface interface {
    Create(network_object *entity.NetworkObject) error

    FindByID(id uint64) (*entity.NetworkObject, error)

    FindByName(name string) (*entity.NetworkObject, error)

    FindAll(page, limit int, sort string) ([]entity.NetworkObject, error)

    Update(network_object *entity.NetworkObject) error

    Delete(id uint64) error
}

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)

}
