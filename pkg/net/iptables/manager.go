package iptables

import (
	"sync"

	"github.com/juju/errors"
	"github.com/sxllwx/vulcanus/pkg/host"
)

// the remote or local iptables-mananger
type Manager struct {
	lock sync.Mutex
	// the shell env for spec host
	host host.Interface
	am   *ArgsManager
}

func (m *Manager) Close() error {
	return m.host.Close()
}

func NewManager(h host.Interface) *Manager {
	return &Manager{
		host: h,
		am:   &ArgsManager{},
	}
}

func (m *Manager) CreateChainForTable(table string, chain string) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	rootCommand := IptablesCommand
	args := m.am.NewChain(table, chain)

	_, err := m.host.Execute(rootCommand, args...)
	if err != nil {
		return errors.Annotatef(err, "iptables create chain (%s) in table (%s)", chain, table)
	}
	return nil
}

func (m *Manager) AppendChainToParentChain(table string, parent string, chain string, comment string) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	rootCommand := IptablesCommand
	args := m.am.AppendChainToParent(table, parent, chain, comment)

	_, err := m.host.Execute(rootCommand, args...)
	if err != nil {
		return errors.Annotatef(err, "iptables create chain (%s) in table (%s)", chain, table)
	}

	return nil
}

func (m *Manager) AppendDNATRuleToChain(
	chain string,
	protocol string,
	// dport can be range eg 1000:2000
	dport string,
	// toDestination can be range eg 192.168.240.101:1000-2000
	toDestination string,
	comment string,
) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	rootCommand := IptablesCommand
	args := m.am.AppendDNATRuleToChain(
		NATTableName,
		chain,
		protocol,
		dport,
		toDestination,
		comment,
	)

	_, err := m.host.Execute(rootCommand, args...)

	if err != nil {
		return errors.Annotatef(err,
			"iptables create dnat rule (protocol %s dport %s to-destination %s comment %s) to chain %s",
			protocol,
			dport,
			toDestination,
			comment,
			chain,
		)
	}

	return nil
}

func (m *Manager) AppendSNATRuleToChain(
	chain string,
	d string,
	comment string,
) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	rootCommand := IptablesCommand
	args := m.am.AppendMASQUERADERuleToChain(
		NATTableName,
		chain,
		d,
		comment,
	)

	_, err := m.host.Execute(rootCommand, args...)
	if err != nil {
		return errors.Annotatef(err,
			"iptables create snat rule (destination %s comment %s) to chain %s",
			d,
			comment,
			chain,
		)
	}

	return nil
}

func (m *Manager) DeleteChain(table string, parent string, chain string, comment string) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	rootCommand := IptablesCommand
	args := m.am.RemoveChainFromParent(NATTableName, parent, chain, comment)

	_, err := m.host.Execute(rootCommand, args...)
	if err != nil {
		return errors.Annotatef(err, "iptables remove chain (%s) from parent (%s) in table (%s)", chain, parent, table)
	}

	args = m.am.FlushChain(NATTableName, chain)
	_, err = m.host.Execute(rootCommand, args...)
	if err != nil {
		return errors.Annotatef(err, "iptables flush chain (%s)  in table (%s)", chain, table)
	}

	args = m.am.DeleteChain(NATTableName, chain)
	_, err = m.host.Execute(rootCommand, args...)
	if err != nil {
		return errors.Annotatef(err, "iptables delete chain (%s)  in table (%s)", chain, table)
	}
	return nil
}

func (m *Manager) checkRule(table string, chain string, args ...string) error {

	rootCommand := IptablesCommand
	realArgs := m.am.CheckRule(table, chain, args...)
	_, err := m.host.Execute(rootCommand, realArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Reset() error {

	m.lock.Lock()
	defer m.lock.Unlock()

	rootCommand := IptablesCommand
	tables := []string{NATTableName, FilterTableName}

	for _, table := range tables {

		args := m.am.FlushAllChain(table)
		_, err := m.host.Execute(rootCommand, args...)
		if err != nil {
			return errors.Annotatef(err, "iptables flush all chain  in table (%s)", table)
		}

		args = m.am.ZeroAllChain(table)
		_, err = m.host.Execute(rootCommand, args...)
		if err != nil {
			return errors.Annotatef(err, "iptables zero all chain  in table (%s)", table)
		}

		args = m.am.ZeroAllChain(table)
		_, err = m.host.Execute(rootCommand, args...)
		if err != nil {
			return errors.Annotatef(err, "iptables zero all chain  in table (%s)", table)
		}
	}

	return nil
}
