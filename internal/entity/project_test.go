package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewProject(t *testing.T) {
	tests := []struct {
		name     string
		tenantID uint64
		err        error
	}{
		{"My project", 1, nil},
		{"", 1, entity.ErrNameRequired},
		{"My project", 0, entity.ErrTenantRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project, err := entity.NewProject(tt.name, tt.tenantID)
			if tt.err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, project)
			} else {
				assert.Equal(t, tt.err, err)
			}
		})
	}
}


func TestNewProjectWithChain(t *testing.T) {
	
	chain, _ := entity.NewChain("INPUT", "filter","accept",1,1,1)
	project, _ := entity.NewProject("My project", 1)
	project.AddChain(chain)
	assert.Equal(t, len(project.Chains), 1)
}