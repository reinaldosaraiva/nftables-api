package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewNetworkObject(t *testing.T) {
    tests := []struct {
        name    string
        address string
        err     error
    }{
        {"Local Network", "192.168.1.0/24", nil},
        {"Main Server", "10.0.0.1", nil},
        {"", "192.168.1.0/24", entity.ErrNameRequired},
        {"Local Network", "", entity.ErrAddressRequired},
        // Mais casos de teste podem ser adicionados aqui
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            networkObject, err := entity.NewNetworkObject(tt.name, tt.address)

            if tt.err == nil {
                assert.Nil(t, err)
                assert.NotNil(t, networkObject)
                assert.Equal(t, tt.name, networkObject.Name)
                assert.Equal(t, tt.address, networkObject.Address)
            } else {
                assert.Equal(t, tt.err, err)
            }
        })
    }
}
