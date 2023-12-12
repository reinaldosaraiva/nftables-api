package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewChain(t *testing.T) {
	tests := []struct {
		name     string
		description string
		type_name       string
		state      string
		projectID  uint64
		err        error
	}{
		{"My chain", "This is a chain", "Filtering", "Active", 1, nil},
		{"", "This is a chain", "Filtering", "Active", 1, entity.ErrNameRequired},
		{"My chain", "", "Filtering", "Active", 1, entity.ErrDescriptionRequired},
		{"My chain", "This is a chain", "", "Active", 1, entity.ErrTypeRequired},
		{"My chain", "This is a chain", "Filtering", "", 1, entity.ErrStateRequired},
		{"My chain", "This is a chain", "Filtering", "Active", 0, entity.ErrProjectRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, err := entity.NewChain(tt.name, tt.description, tt.type_name, tt.state, tt.projectID)
			if tt.err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, chain)
			} else {
				assert.Equal(t, err, tt.err)
			}
		})
	}
}
