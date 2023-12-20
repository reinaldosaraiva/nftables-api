package database

import (
	"fmt"
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDatabaseForChain(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }

    err = db.AutoMigrate(&entity.Project{}, &entity.Table{}, &entity.Chain{}, &entity.Rule{})
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }

    project := &entity.Project{Name: "Project 1", TenantID: 1}
    db.Create(project)
    table := &entity.Table{Name: "mangle", Description: "Mangle table",Type: "mangle", }
    db.Create(table)

    return db
}

func TestCreateChain(t *testing.T) {
    tests := []struct {
        name      string
        type_name string
        priority  int
        policy    string
        tableID   uint64
        err       error
    }{
        {"PREROUTING", "filter", 1, "ACCEPT", 1, nil},
        // {"INPUT", "filter", 0, "ACCEPT", 1, nil},
        // {"input", "filter", 0, "DROP", 4, nil},
        // {"output", "filter", 0, "ACCEPT", 4, nil},
    }


    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            chain := &entity.Chain{
                Name:      tt.name,
                Type:      tt.type_name,
                Priority:  tt.priority,
                Policy:    tt.policy,
                TableID:   tt.tableID,
                ProjectID: 1,
            }
            err := chainDB.Create(chain)
            if tt.err != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.err, err)
            } else {
                assert.NoError(t, err)
                assert.NotZero(t, chain.ID)
            }
        })
    }
}


func TestFindChainByID(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    chain := &entity.Chain{Name: "INPUT", Type: "filter",Priority: 1 , ProjectID: 1, TableID: 1}
    err := chainDB.Create(chain)
    assert.NoError(t, err)

    foundChain, err := chainDB.FindByID(uint64(chain.ID))
    assert.NoError(t, err)
    assert.NotNil(t, foundChain)
    assert.Equal(t, "INPUT", foundChain.Name)
}

func TestFindChainByName(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    chain := &entity.Chain{Name: "INPUT", Type: "filter",Priority: 1 , ProjectID: 1, TableID: 1}
    err := chainDB.Create(chain)
    assert.NoError(t, err)

    foundChain, err := chainDB.FindByName("INPUT")
    assert.NoError(t, err)
    assert.NotNil(t, foundChain)
    assert.Equal(t, "INPUT", foundChain.Name)
}

func TestUpdateChain(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    chain := &entity.Chain{Name: "INPUT", Type: "filter",Priority: 1 , ProjectID: 1, TableID: 1}
    err := chainDB.Create(chain)
    assert.NoError(t, err)

    chain.Name = "Updated Chain"
    err = chainDB.Update(chain)
    assert.NoError(t, err)

    updatedChain, err := chainDB.FindByID(uint64(chain.ID))
    assert.NoError(t, err)
    assert.Equal(t, "Updated Chain", updatedChain.Name)
}

func TestFindAllChains(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    for i := 0; i < 10; i++ {
        chain := &entity.Chain{Name: fmt.Sprintf("Chain %d", i),  Type: "filter",Priority: 1, ProjectID: 1, TableID: 1}
        err := chainDB.Create(chain)
        assert.NoError(t, err)
    }

    chains, err := chainDB.FindAll(1, 5, "asc")
    assert.NoError(t, err)
    assert.Len(t, chains, 5)
}

func TestDeleteChain(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    chain := &entity.Chain{Name: "INPUT", Type: "filter",Priority: 1 , ProjectID: 1, TableID: 1}
    err := chainDB.Create(chain)
    assert.NoError(t, err)

    err = chainDB.Delete(uint64(chain.ID))
    assert.NoError(t, err)

    _, err = chainDB.FindByID(uint64(chain.ID))
    assert.Error(t, err)
}
