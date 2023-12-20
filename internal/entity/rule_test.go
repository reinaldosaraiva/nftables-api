package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewRule(t *testing.T) {
	tests := []struct {
		chainID   uint64
		protocol  string
		port      int
		action    string
		wantErr   bool
		errMsg    string
	}{
		{1, "TCP", 80, "ACCEPT", false, ""},
		{0, "TCP", 80, "ACCEPT", true, "Chain is required"},
		{1, "TCP", 0, "ACCEPT", true, "Port is required"},
		// Outros casos de teste...
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			rule, err := entity.NewRule(tt.chainID, tt.protocol, tt.port, tt.action)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rule)
			}
		})
	}
}

func TestNewRuleWithRelations(t *testing.T) {

	chain := entity.Chain{Name: "Test Chain", Type: "filter", Priority: 1, Policy: "accept"}
	service := entity.Service{Name: "Test Service"}
	networkObject := entity.NetworkObject{Name: "Test Network Object"}

	rule, err := entity.NewRule(1, "TCP", 80, "ACCEPT")
	assert.NoError(t, err)
	assert.NotNil(t, rule)


	rule.Chain = chain
	rule.ServiceRules = []entity.Service{service}
	rule.NetworkObjectRules = []entity.NetworkObject{networkObject}


	assert.Equal(t, "Test Chain", rule.Chain.Name)
	assert.Equal(t, "Test Service", rule.ServiceRules[0].Name)
	assert.Equal(t, "Test Network Object", rule.NetworkObjectRules[0].Name)
}