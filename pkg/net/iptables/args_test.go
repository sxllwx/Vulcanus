package iptables

import (
	"testing"
)

func TestNewIPTableManager(t *testing.T) {

	m := &ArgsManager{}

	const (
		ChainName = "DESKTOP-SERVICES"
	)
	t.Log(m.NewChain(NATTableName, ChainName))
	t.Log(m.DeleteChain(NATTableName, ChainName))
	t.Log(m.FlushChain(NATTableName, ChainName))
	t.Log(m.AppendChainToParent(NATTableName, PREROUTINGChainName, ChainName, ""))
	t.Log(m.AppendDNATRuleToChain(
		NATTableName,
		ChainName,
		"tcp",
		"3000:4000",
		"192.168.240.101:3000-4000",
		"",
	))
	t.Log(m.AppendDNATRuleToChain(
		NATTableName,
		ChainName,
		"tcp",
		"3000:4000",
		"192.168.240.101:3000-4000",
		"",
	))

	t.Log(m.CheckRule(NATTableName, PREROUTINGChainName, ""))
}
