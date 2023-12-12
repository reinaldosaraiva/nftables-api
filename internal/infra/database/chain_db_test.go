// chain_db_test.go
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

    err = db.AutoMigrate(&entity.Chain{}, &entity.Table{}, &entity.Rule{})
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }

    return db
}

func TestCreateChain(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)
    chain := &entity.Chain{Name: "Chain 1", Type: "SomeType", State: "Active", ProjectID: 1, TableID: 1} 

    err := chainDB.Create(chain)
    assert.NoError(t, err)
    assert.NotZero(t, chain.ID)
}

func TestFindChainByID(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    chain := &entity.Chain{Name: "Chain 1", Type: "SomeType", State: "Active", ProjectID: 1, TableID: 1}
    err := chainDB.Create(chain)
    assert.NoError(t, err)

    foundChain, err := chainDB.FindByID(chain.ID)
    assert.NoError(t, err)
    assert.NotNil(t, foundChain)
    assert.Equal(t, "Chain 1", foundChain.Name)
}

func TestUpdateChain(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    chain := &entity.Chain{Name: "Chain 1", Type: "SomeType", State: "Active", ProjectID: 1, TableID: 1}
    err := chainDB.Create(chain)
    assert.NoError(t, err)

    chain.Name = "Updated Chain"
    err = chainDB.Update(chain)
    assert.NoError(t, err)

    updatedChain, err := chainDB.FindByID(chain.ID)
    assert.NoError(t, err)
    assert.Equal(t, "Updated Chain", updatedChain.Name)
}

func TestFindAllChains(t *testing.T) {
    db := setupDatabaseForChain(t)
    chainDB := NewChainDB(db)

    for i := 0; i < 10; i++ {
        chain := &entity.Chain{Name: fmt.Sprintf("Chain %d", i), Type: "SomeType", State: "Active", ProjectID: 1, TableID: 1}
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

    chain := &entity.Chain{Name: "Chain to Delete", Type: "SomeType", State: "Active", ProjectID: 1, TableID: 1}
    err := chainDB.Create(chain)
    assert.NoError(t, err)

    err = chainDB.Delete(chain.ID)
    assert.NoError(t, err)

    _, err = chainDB.FindByID(chain.ID)
    assert.Error(t, err)
}
