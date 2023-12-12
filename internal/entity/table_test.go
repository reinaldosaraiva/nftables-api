package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewTable(t *testing.T) {
	tests := []struct {
		name     string
		description string
		type_name       string
		state      string
		err        error
	}{
		{"My table", "This is a table", "Database", "Active", nil},
		{"", "This is a table", "Database", "Active", entity.ErrNameRequired},
		{"My table", "", "Database", "Active", entity.ErrDescriptionRequired},
		{"My table", "This is a table", "", "Active", entity.ErrTypeRequired},
		{"My table", "This is a table", "Database", "", entity.ErrStateRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := entity.NewTable(tt.name, tt.description, tt.type_name, tt.state)
			if tt.err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, table)
			} else {
				assert.Equal(t, err, tt.err)
			}
		})
	}
}
