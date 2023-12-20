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

    err = db.AutoMigrate(&entity.Project{}, &entity.Chain{}, &entity.Table{})
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }

    testProject := entity.Project{Name: "Test Project"}
    if err := db.Create(&testProject).Error; err != nil {
        t.Fatalf("Failed to create test project: %v", err)
    }

    return db
}


func TestCreateTable(t *testing.T) {
    tests := []struct {
        name        string
        description string
        type_name   string
        err         error
        
    }{
        {"mangle", "Mangle table", "mangle",  nil},
        {"security", "Security table", "security",  nil},
        {"raw", "Raw table", "raw",  nil},
        {"inet", "INET table (main filter)", "filter",  nil},
		{"inet6", "INET6 table (main filter)", "filter",  nil},
		{"arp", "ARP table (main filter)", "filter",   nil},
		{"bridge", "Bridge table (main filter)", "filter",   nil},
		{"ip", "IP table (main filter)", "filter",  nil},
		{"netdev", "Netdev table (main filter)", "filter",  nil},
		{"route", "Route table (main filter)", "filter", nil},
		{"fib", "FIB table (main filter)", "filter",  nil},
		{"fip", "FIP table (main filter)", "filter",  nil},
		{"inet", "INET table (main nat)", "nat", nil},
    }

    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            table := &entity.Table{
                Name:        tt.name,
                Description: tt.description,
                Type:        tt.type_name,

            } 

            err := tableDB.Create(table)
            if tt.err != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.err, err)
            } else {
                assert.NoError(t, err)
                assert.NotZero(t, table.ID)
            }
        })
    }
}


func TestCreateTableWithChains(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    chains := []entity.Chain{
        {Name: "INPUT", Type: "filter",Priority: 1 , ProjectID: 1, TableID: 1},
        {Name: "OUTPUT", Type: "filter",Priority: 1 , ProjectID: 1, TableID: 1},
    }

    for _, chain := range chains {
        err := db.Create(&chain).Error
        assert.NoError(t, err)
    }

    table := &entity.Table{
        Name: "mangle", Description: "Mangle table",Type: "mangle",
        Chains: chains, 
    }

    err := tableDB.Create(table)
    assert.NoError(t, err)
    assert.NotZero(t, table.ID)

    assert.Len(t, table.Chains, len(chains))
}

func TestFindTableByID(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    table := &entity.Table{Name: "mangle", Description: "Mangle table",Type: "mangle", }
    err := tableDB.Create(table)
    assert.NoError(t, err)

    foundTable, err := tableDB.FindByID(uint64(table.ID))
    assert.NoError(t, err)
    assert.NotNil(t, foundTable)
    assert.Equal(t, "mangle", foundTable.Name)
}

func TestUpdateTable(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    table := &entity.Table{Name: "mangle", Description: "Mangle table",Type: "mangle",}
    err := tableDB.Create(table)
    assert.NoError(t, err)

    table.Name = "Updated Table"
    err = tableDB.Update(table)
    assert.NoError(t, err)

    updatedTable, err := tableDB.FindByID(uint64(table.ID))
    assert.NoError(t, err)
    assert.Equal(t, "Updated Table", updatedTable.Name)
}

func TestFindTableByName(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    table := &entity.Table{Name: "mangle", Description: "Mangle table",Type: "mangle", }
    err := tableDB.Create(table)
    assert.NoError(t, err)

    foundTable, err := tableDB.FindByName("mangle")
    assert.NoError(t, err)
    assert.NotNil(t, foundTable)
    assert.Equal(t, "mangle", foundTable.Name)
}
func TestFindAllTables(t *testing.T) {
    db := setupDatabaseForTable(t)
    tableDB := NewTableDB(db)

    for i := 0; i < 10; i++ {
        table := &entity.Table{Name: fmt.Sprintf("Table %d", i), Type: "SomeType", Description: "Some table",}
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

    table := &entity.Table{Name: "mangle", Description: "Mangle table",Type: "mangle",  }
    err := tableDB.Create(table)
    assert.NoError(t, err)

    err = tableDB.Delete(uint64(table.ID))
    assert.NoError(t, err)

    _, err = tableDB.FindByID(uint64(table.ID))
    assert.Error(t, err)
}

