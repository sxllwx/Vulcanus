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
	BlackListChainName   = "BLACK-LIST"
	PREROUTINGChainName  = "PREROUTING"
	POSTROUTINGChainName = "POSTROUTING"
	INPUTChainName       = "INPUT"
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

func (t *IPTableManager) NewChain(table string, nativeChain string, chain string) ([]byte, error) {

	out, err := t.Runner(
		iptablesCommand,
		"-t", table,
		"-N", chain,
	)
	if err != nil {
		return out, errors.Annotate(err, "execute new chain command")
	}
	out, err = t.Runner(
		iptablesCommand,
		"-t", table,
		"-I", nativeChain,
		"-j", chain,
	)
	if err != nil {
		return out, errors.Annotate(err, "execute insert new chain to native chain command")
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
		"--to-destination", fmt.Sprintf("%s:%s", rule.ToHost, rule.ToPort),
	)

	if err != nil {
		return out, errors.Annotate(err, "execute dnat command")
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
		return out, errors.Annotate(err, "execute snat command")
	}
	return out, nil
}
