package iptables

import (
	"github.com/sxllwx/vulcanus/pkg/host/localhost"
	"testing"
)

func localIPTablesManager() *IPTableManager {

	local := localhost.NewClient(&localhost.Config{})
	return NewIPTableManager(local.Execute)
}

func TestNewIPTableManager(t *testing.T) {

	m := localIPTablesManager()

	out, err := m.NewChain(FilterTableName, "INPUT", BlackListChainName)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", out)

}
