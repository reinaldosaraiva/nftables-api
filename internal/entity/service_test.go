package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
    tests := []struct {
        name    string
        port    int
        err     error
    }{
        {"HTTP", 80, nil},
        {"HTTPS", 443, nil},
        {"SSH", 22, nil},
        {"FTP", 21, nil},
        {"SMTP", 25, nil},
        {"", 80, entity.ErrNameRequired},
        {"HTTP", 0, entity.ErrPortRequired},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service, err := entity.NewService(tt.name, tt.port)
            if tt.err == nil {
                assert.Nil(t, err)
                assert.NotNil(t, service)
                assert.Equal(t, tt.name, service.Name)
                assert.Equal(t, tt.port, service.Port)
            } else {
                assert.Equal(t, err, tt.err)
            }
        })
    }
}
