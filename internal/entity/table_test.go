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
        priority    int
        comment     string
        err         error
    }{
        {"mangle", "Mangle table", "mangle", 100, "mangle", nil},
        {"security", "Security table", "security", 150, "security", nil},
        {"raw", "Raw table", "raw", 200, "raw", nil},
        {"inet", "INET table (main filter)", "filter", 300, "filter", nil},
		{"inet6", "INET6 table (main filter)", "filter", 300, "filter", nil},
		{"arp", "ARP table (main filter)", "filter", 300, "filter", nil},
		{"bridge", "Bridge table (main filter)", "filter", 300, "filter", nil},
		{"ip", "IP table (main filter)", "filter", 300, "filter", nil},
		{"netdev", "Netdev table (main filter)", "filter", 300, "filter", nil},
		{"route", "Route table (main filter)", "filter", 300, "filter", nil},
		{"fib", "FIB table (main filter)", "filter", 300, "filter", nil},
		{"fip", "FIP table (main filter)", "filter", 300, "filter", nil},
		{"inet", "INET table (main nat)", "nat", 300, "filter", nil},
		{"", "This is a table", "Database", 0,  "filter", entity.ErrNameRequired},
        {"My table", "", "Database", 0,  "filter", entity.ErrDescriptionRequired},
        {"My table", "This is a table", "", 0, "filter", entity.ErrTypeRequired},
		{"My table", "This is a table", "Database", 0, "", entity.ErrCommentRequired},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            table, err := entity.NewTable(tt.name, tt.description, tt.type_name, tt.priority, tt.comment)
            if tt.err == nil {
                assert.Nil(t, err)
                assert.NotNil(t, table)
            } else {
                assert.Equal(t, err, tt.err)
            }
        })
    }
}
