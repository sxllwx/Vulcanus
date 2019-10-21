package iptables

import (
	"fmt"

	"github.com/juju/errors"
)

// root command
const (
	iptablesCommand = "iptables"
)

// tables name
const (
	NATTableName    = "nat"
	FilterTableName = "filter"
)

// default chain
const (
	BlackListChainName = "BLACK-LIST"
)

type IPTableManager struct {
	Runner func(cmd string, args ...string) ([]byte, error)
}

func NewIPTableManager(runner func(cmd string, args ...string) ([]byte, error)) *IPTableManager {

	m := &IPTableManager{
		Runner: runner,
	}
	return m
}

func (t *IPTableManager) NewChain(table string, chain string) ([]byte, error) {

	out, err := t.Runner(
		iptablesCommand,
		"-t", table,
		"-N", chain,
	)
	if err != nil {
		return out, errors.Annotatef(err, "new chain %s in table %s, the sys-out %s", chain, table, out)
	}
	return out, nil
}

func (t *IPTableManager) DNAT(chain string, rule DNATRule) ([]byte, error) {

	out, err := t.Runner(
		iptablesCommand,
		"-t", NATTableName,
		"-I", chain,
		"-p", "tcp",
		"--dport", rule.DPort,
		"-j", "DNAT",
		"-to-destination", fmt.Sprintf("%s:%s", rule.ToHost, rule.ToPort),
	)

	if err != nil {
		return out, errors.Annotatef(err, "new dnat rule %#v in chain %s, the sys-out %s", rule, chain, out)
	}
	return out, nil

}

func (t *IPTableManager) SNAT(chain string, rule SNATRule) ([]byte, error) {

	out, err := t.Runner(
		iptablesCommand,
		"-t", NATTableName,
		"-I", chain,
		"-p", "tcp",
		"-d", rule.DHost,
		"-j", "MASQUERADE",
	)

	if err != nil {
		return out, errors.Annotatef(err, "new snat rule %#v in chain %s, the sys-out %s", rule, chain, out)
	}
	return out, nil
}
