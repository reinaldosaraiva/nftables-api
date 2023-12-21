package database

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDatabaseForRule(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }

    err = db.AutoMigrate(&entity.Chain{}, &entity.Rule{},&entity.Project{}, &entity.Table{},&entity.Tenant{},&entity.Service{},&entity.NetworkObject{},)
    if err != nil {
        t.Fatalf("Failed to migrate database: %v", err)
    }


    chain := &entity.Chain{Name: "INPUT", Type: "filter",Priority: 1 , ProjectID: 1, TableID: 1}
    db.Create(chain)

    return db
}

func TestCreateRule(t *testing.T) {
    db := setupDatabaseForRule(t)
    ruleDB := NewRuleDB(db)

    rule := &entity.Rule{
        ChainID:  1,
        Protocol: "TCP",
        Port:     80,
        Action:   "ACCEPT",
    }

    err := ruleDB.Create(rule)
    assert.NoError(t, err)
    assert.NotZero(t, rule.ID)
}

func TestFindRuleByID(t *testing.T) {
    db := setupDatabaseForRule(t)
    ruleDB := NewRuleDB(db)

    rule := &entity.Rule{
        ChainID:  1,
        Protocol: "TCP",
        Port:     80,
        Action:   "ACCEPT",
    }
    err := ruleDB.Create(rule)
    assert.NoError(t, err)

    foundRule, err := ruleDB.FindByID(uint64(rule.ID))
    assert.NoError(t, err)
    assert.NotNil(t, foundRule)
    assert.Equal(t, 80, foundRule.Port)
}

func TestUpdateRule(t *testing.T) {
    db := setupDatabaseForRule(t)
    ruleDB := NewRuleDB(db)

    rule := &entity.Rule{
        ChainID:  1,
        Protocol: "TCP",
        Port:     80,
        Action:   "ACCEPT",
    }
    err := ruleDB.Create(rule)
    assert.NoError(t, err)

    rule.Port = 443
    err = ruleDB.Update(rule)
    assert.NoError(t, err)

    updatedRule, err := ruleDB.FindByID(uint64(rule.ID))
    assert.NoError(t, err)
    assert.Equal(t, 443, updatedRule.Port)
}

func TestFindAllRules(t *testing.T) {
    db := setupDatabaseForRule(t)
    ruleDB := NewRuleDB(db)

    for i := 0; i < 10; i++ {
        rule := &entity.Rule{
            ChainID:  1,
            Protocol: "TCP",
            Port:     80 + i,
            Action:   "ACCEPT",
        }
        err := ruleDB.Create(rule)
        assert.NoError(t, err)
    }

    rules, err := ruleDB.FindAll(1, 5, "asc")
    assert.NoError(t, err)
    assert.Len(t, rules, 5)
}

func TestDeleteRule(t *testing.T) {
    db := setupDatabaseForRule(t)
    ruleDB := NewRuleDB(db)

    rule := &entity.Rule{
        ChainID:  1,
        Protocol: "TCP",
        Port:     80,
        Action:   "ACCEPT",
    }
    err := ruleDB.Create(rule)
    assert.NoError(t, err)

    err = ruleDB.Delete(uint64(rule.ID))
    assert.NoError(t, err)

    _, err = ruleDB.FindByID(uint64(rule.ID))
    assert.Error(t, err)
}


func TestCreateRuleWithRelations(t *testing.T) {
    db := setupDatabaseForRule(t)
    ruleDB := NewRuleDB(db)

    
    chain := &entity.Chain{ Name: "Chain 1", Type: "Type",  ProjectID: 1, TableID: 1}
    db.Create(chain)

    service := &entity.Service{Name: "HTTP", Port: 80}
    db.Create(service)

    networkObject := &entity.NetworkObject{Name: "NetworkObject 1", Address: "192.168.0.0/24"}
    db.Create(networkObject)

    
    rule := &entity.Rule{
        ChainID:    uint64(chain.ID),
        Port:       80,
        Protocol:   "TCP",
        Action:     "ACCEPT",
        ChainRules: []entity.Chain{*chain},
        ServiceRules: []entity.Service{*service},
        NetworkObjectRules: []entity.NetworkObject{*networkObject},
    }

    err := ruleDB.Create(rule)
    assert.NoError(t, err)

    var foundRule entity.Rule
    err = db.Preload("ChainRules").Preload("ServiceRules").Preload("NetworkObjectRules").First(&foundRule, rule.ID).Error
    assert.NoError(t, err)
    assert.NotNil(t, foundRule.ChainRules)
    assert.NotNil(t, foundRule.ServiceRules)
    assert.NotNil(t, foundRule.NetworkObjectRules)


    assert.Equal(t, chain.Name, foundRule.ChainRules[0].Name)
    assert.Equal(t, service.Name, foundRule.ServiceRules[0].Name)
    assert.Equal(t, networkObject.Name, foundRule.NetworkObjectRules[0].Name)
}
