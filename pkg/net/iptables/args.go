package iptables

import (
	"fmt"
)

// root command
const (
	IptablesCommand     = "iptables"
	IptablesSaveCommand = "iptables-save"
)

// tables name
const (
	NATTableName    = "nat"
	FilterTableName = "filter"
)

// default chain
const (
	PREROUTINGChainName  = "PREROUTING"
	POSTROUTINGChainName = "POSTROUTING"
	INPUTChainName       = "INPUT"
	OUTPUTChainName      = "OUTPUT"
)

// default Action
const (
	ACCEPT     = "ACCEPT"
	DROP       = "DROP"
	DNAT       = "DNAT"
	SNAT       = "SNAT"
	MASQUERADE = "MASQUERADE"
)

// default iptables comment
const (
	DefaultComment = "generated by vulcanus"
)

type ArgsManager struct{}

// NewChain
// new a chain in spec table
// eg: KUBE-SERVICES
// eg: KUBE-SVC-XXXX
// eg: KUBE-SEP-XXXX
func (m *ArgsManager) NewChain(table string, chain string) []string {

	return append(
		[]string{},
		"-t", table,
		"-N", chain,
	)
}

// Delete chain
// delete a chain from table
func (m *ArgsManager) DeleteChain(table string, chain string) []string {

	return append(
		[]string{},
		"-t", table,
		"-X", chain,
	)
}

// Flush chain
// flush a chain in spec table
func (m *ArgsManager) FlushChain(table string, chain string) []string {

	return append(
		[]string{},
		"-t", table,
		"-F", chain,
	)
}

// InsertToParent
// insert a chain to parent chain in spec table
// with comment
func (m *ArgsManager) AppendChainToParent(
	table string,
	parent string,
	chain string,
	comment string,
) []string {

	if len(comment) == 0 {
		comment = DefaultComment
	}

	return append(
		[]string{},
		"-t", table,
		"-A", parent,
		"-m", "comment",
		"--comment", m.wrapComment(comment),
		"-j", chain,
	)
}

func (m *ArgsManager) wrapComment(comment string) string {
	return fmt.Sprintf("\"%s\"", comment)
}

// RemoveChainFromParent
// remove the chain from parent chain in spec table
func (m *ArgsManager) RemoveChainFromParent(
	table string,
	parent string,
	chain string,
	comment string,
) []string {

	if len(comment) == 0 {
		comment = DefaultComment
	}
	return append(
		[]string{},
		"-t", table,
		"-D", parent,
		"-m", "comment",
		"--comment", m.wrapComment(comment),
		"-j", chain,
	)
}

// CheckRule
// check the iptables  rule exist
func (m *ArgsManager) CheckRule(table string, chain string, args ...string) []string {

	return append(
		[]string{
			"-t", table,
			"-C", chain,
		},
		args...,
	)
}

// AppendDNATRuleToChain
// append a dnat rule for a chain in spec table
func (m *ArgsManager) AppendDNATRuleToChain(
	table string,
	chain string,
	protocol string,
	// dport can be range eg 1000:2000
	dport string,
	// toDestination can be range eg 192.168.240.101:1000-2000
	toDestination string,
	comment string,
) []string {

	if len(comment) == 0 {
		comment = DefaultComment
	}

	return append(
		[]string{},
		"-t", table,
		"-A", chain,
		"-p", protocol,
		"--dport", dport,
		"-m", "comment",
		"--comment", m.wrapComment(comment),
		"-j", DNAT,
		"--to-destination", toDestination,
	)
}

// AppendMASQUERADERuleToChain
// append a snat rule for a chain in spec table
func (m *ArgsManager) AppendMASQUERADERuleToChain(
	table string,
	chain string,
	d string,
	comment string,
) []string {

	if len(comment) == 0 {
		comment = DefaultComment
	}

	return append(
		[]string{},
		"-t", table,
		"-A", chain,
		"-d", d,
		"-m", "comment",
		"--comment", m.wrapComment(comment),
		"-j", MASQUERADE,
	)
}

// FlushAllChain
// Flush all chain in spec table
func (m *ArgsManager) FlushAllChain(table string) []string {

	return append(
		[]string{},
		"-t", table,
		"-F",
	)
}

// DeleteAllChain
// Delete all chain in spec table
func (m *ArgsManager) DeleteAllChain(table string) []string {

	return append(
		[]string{},
		"-t", table,
		"-X",
	)
}

// ZeroAllChain
// Zero all chain in spec table
func (m *ArgsManager) ZeroAllChain(table string) []string {

	return append(
		[]string{},
		"-t", table,
		"-Z",
	)
}
