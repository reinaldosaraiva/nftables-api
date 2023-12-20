package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewChain(t *testing.T) {
    tests := []struct {
        name        string
        type_name   string
        priority    int
        policy      string
        projectID   uint64
        tableID     uint64
        err         error
    }{
        {"My chain", "Filtering", 100, "ACCEPT", 1, 1, nil},
        {"", "Filtering", 100, "ACCEPT", 1, 1, entity.ErrNameRequired},
        {"My chain", "", 100, "ACCEPT", 1, 1, entity.ErrTypeRequired},
        {"My chain", "Filtering", 100, "", 1, 1, entity.ErrPolicyRequired},
        {"My chain", "Filtering", 100, "ACCEPT", 0, 1, entity.ErrProjectRequired},
        {"My chain", "Filtering", 100, "ACCEPT", 1, 0, entity.ErrTableRequired},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            chain, err := entity.NewChain(tt.name, tt.type_name,  tt.policy,tt.priority, tt.projectID, tt.tableID)
            if tt.err == nil {
                assert.Nil(t, err)
                assert.NotNil(t, chain)
            } else {
                assert.Equal(t, err, tt.err)
            }
        })
    }
}
