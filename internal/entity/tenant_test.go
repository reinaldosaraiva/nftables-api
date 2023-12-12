package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewTenant(t *testing.T) {
	tests := []struct {
		name     string
		err        error
	}{
		{"My tenant", nil},
		{"", entity.ErrNameRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenant, err := entity.NewTenant(tt.name)
			if tt.err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, tenant)
			} else {
				assert.Equal(t, err, tt.err)
			}
		})
	}
}
