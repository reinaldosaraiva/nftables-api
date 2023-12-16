// table_db_test.go
package database

import (
	"fmt"
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDatabaseForTable(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }

    // Adicione a migração para todas as entidades relacionadas
    err = db.AutoMigrate(&entity.Table{}, &entity.Project{}, &entity.Tenant{}, &entity.Chain{})
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }

    return db
}

func TestCreateTable(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)
    table := &entity.Table{Name: "Table 1", Type: "SomeType", State: "Active" } // Suponha que TenantID 1 exista

    err := tableDB.Create(table)
    assert.NoError(t, err)
    assert.NotZero(t, table.ID)
}

func TestFindTableByID(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    table := &entity.Table{Name: "Table 1", Type: "SomeType", State: "Active"}
    err := tableDB.Create(table)
    assert.NoError(t, err)

    foundTable, err := tableDB.FindByID(uint64(table.ID))
    assert.NoError(t, err)
    assert.NotNil(t, foundTable)
    assert.Equal(t, "Table 1", foundTable.Name)
}

func TestUpdateTable(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    table := &entity.Table{Name: "Table 1", Type: "SomeType", State: "Active"}
    err := tableDB.Create(table)
    assert.NoError(t, err)

    table.Name = "Updated Table"
    err = tableDB.Update(table)
    assert.NoError(t, err)

    updatedTable, err := tableDB.FindByID(uint64(table.ID))
    assert.NoError(t, err)
    assert.Equal(t, "Updated Table", updatedTable.Name)
}

func TestFindAllTables(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    for i := 0; i < 10; i++ {
        table := &entity.Table{Name: fmt.Sprintf("Table %d", i), Type: "SomeType", State: "Active"}
        err := tableDB.Create(table)
        assert.NoError(t, err)
    }

    tables, err := tableDB.FindAll(1, 5, "asc")
    assert.NoError(t, err)
    assert.Len(t, tables, 5)
}

func TestDeleteTable(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    table := &entity.Table{Name: "Table to Delete", Type: "SomeType", State: "Active"}
    err := tableDB.Create(table)
    assert.NoError(t, err)

    err = tableDB.Delete(uint64(table.ID))
    assert.NoError(t, err)

    _, err = tableDB.FindByID(uint64(table.ID))
    assert.Error(t, err)
}
