package database

import (
	"fmt"
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDatabaseForTenant(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }


    err = db.AutoMigrate(&entity.Tenant{}, &entity.Project{},&entity.Table{},&entity.Chain{})
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }

    return db
}

func TestCreateTenant(t *testing.T) {
	db := setupDatabaseForTenant(t)
	tenantDB := NewTenantDB(db)
	tenant, err := entity.NewTenant("Tenant 1")
	assert.NoError(t, err)

	err = tenantDB.Create(tenant)
	assert.NoError(t, err)
	assert.NotZero(t, tenant.ID)
}

func TestFindTenantByID(t *testing.T) {
	db := setupDatabaseForTenant(t)
	tenantDB := NewTenantDB(db)

	tenant, err := entity.NewTenant("Tenant 1")
	assert.NoError(t, err)
	err = tenantDB.Create(tenant)
	assert.NoError(t, err)

	foundTenant, err := tenantDB.FindByID(uint64(tenant.ID))
	assert.NoError(t, err)
	assert.NotNil(t, foundTenant)
	assert.Equal(t, "Tenant 1", foundTenant.Name)
}

func TestFindTenantByName(t *testing.T) {
    db := setupDatabaseForTenant(t)
    tenantDB := NewTenantDB(db)

    tenant, err := entity.NewTenant("Tenant 1")
    assert.NoError(t, err)
    err = tenantDB.Create(tenant)
    assert.NoError(t, err)

    foundTenant, err := tenantDB.FindByName("Tenant 1")
    assert.NoError(t, err)
    assert.NotNil(t, foundTenant)
    assert.Equal(t, "Tenant 1", foundTenant.Name)
}

func TestDeleteTenantCascade(t *testing.T) {
    db := setupDatabaseForTenant(t)
    tenantDB := NewTenantDB(db)
    projectDB := NewProjectDB(db)
    tableDB := NewTableDB(db) 
    chainDB := NewChainDB(db)

    tenant := &entity.Tenant{Name: "Test Tenant"}
    err := tenantDB.Create(tenant)
    assert.NoError(t, err)

    project := &entity.Project{Name: "Project for Tenant", TenantID: uint64(tenant.ID)}
    err = projectDB.Create(project)
    assert.NoError(t, err)

    table := &entity.Table{Name: "mangle", Description: "Mangle table",Type: "mangle", }
    err = tableDB.Create(table) 
    assert.NoError(t, err)

    for i := 0; i < 3; i++ {
        chain := &entity.Chain{Name: fmt.Sprintf("Chain %d", i), ProjectID: uint64(project.ID), TableID: uint64(table.ID)}
        err = chainDB.Create(chain)
        assert.NoError(t, err)
    }

    err = tenantDB.Delete(uint64(tenant.ID))
    assert.NoError(t, err)

    _, err = tenantDB.FindByID(uint64(tenant.ID))
    assert.Error(t, err)

    _, err = projectDB.FindByID(uint64(project.ID))
    assert.Error(t, err)

    // chains, err := chainDB.FindAll(1, 10, "asc")
    // assert.NoError(t, err)
    // for _, c := range chains {
    //     assert.NotEqual(t, c.ProjectID, project.ID)
    // }
}
