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

    // Adicione a migração para todas as entidades relacionadas
    err = db.AutoMigrate(&entity.Tenant{}, &entity.Project{},&entity.Chain{})
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

func TestUpdateTenant(t *testing.T) {
	db := setupDatabaseForTenant(t)
	tenantDB := NewTenantDB(db)

	tenant, err := entity.NewTenant("Tenant 1")
	assert.NoError(t, err)
	err = tenantDB.Create(tenant)
	assert.NoError(t, err)

	tenant.Name = "Updated Tenant"
	err = tenantDB.Update(tenant)
	assert.NoError(t, err)

	updatedTenant, err := tenantDB.FindByID(uint64(tenant.ID))
	assert.NoError(t, err)
	assert.Equal(t, "Updated Tenant", updatedTenant.Name)
}

func TestFindAllTenants(t *testing.T) {
	db := setupDatabaseForTenant(t)
	tenantDB := NewTenantDB(db)

	for i := 0; i < 10; i++ {
		tenant, err := entity.NewTenant(fmt.Sprintf("Tenant %d", i))
		assert.NoError(t, err)
		err = tenantDB.Create(tenant)
		assert.NoError(t, err)
	}

	tenants, err := tenantDB.FindAll(1, 5, "asc")
	assert.NoError(t, err)
	assert.Len(t, tenants, 5)
}

func TestDeleteTenant(t *testing.T) {
	db := setupDatabaseForTenant(t)
	tenantDB := NewTenantDB(db)

	tenant, err := entity.NewTenant("Tenant to Delete")
	assert.NoError(t, err)
	err = tenantDB.Create(tenant)
	assert.NoError(t, err)

	err = tenantDB.Delete(uint64(tenant.ID))
	assert.NoError(t, err)

	_, err = tenantDB.FindByID(uint64(tenant.ID))
	assert.Error(t, err)
}
