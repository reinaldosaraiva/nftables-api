package entity_test

import (
	"testing"

	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewTable(t *testing.T) {
    tests := []struct {
        name        string
        description string
        type_name   string

        err         error
    }{
        {"mangle", "Mangle table", "mangle", nil},
        {"security", "Security table", "security",  nil},
        {"raw", "Raw table", "raw", nil},
        {"inet", "INET table (main filter)", "filter", nil},
		{"inet6", "INET6 table (main filter)", "filter", nil},
		{"arp", "ARP table (main filter)", "filter", nil},
		{"bridge", "Bridge table (main filter)", "filter", nil},
		{"ip", "IP table (main filter)", "filter", nil},
		{"netdev", "Netdev table (main filter)", "filter", nil},
		{"route", "Route table (main filter)", "filter",  nil},
		{"fib", "FIB table (main filter)", "filter",  nil},
		{"fip", "FIP table (main filter)", "filter",  nil},
		{"inet", "INET table (main nat)", "nat", nil},
		{"", "This is a table", "Database",  entity.ErrNameRequired},
        {"My table", "", "Database",  entity.ErrDescriptionRequired},
		{"My table", "This is a table", "Database", entity.ErrCommentRequired},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            table, err := entity.NewTable(tt.name, tt.description, tt.type_name, )
            if tt.err == nil {
                assert.Nil(t, err)
                assert.NotNil(t, table)
            } else {
                assert.Equal(t, err, tt.err)
            }
        })
    }
}
