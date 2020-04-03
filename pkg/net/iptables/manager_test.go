package iptables

import (
	"log"
	"os"
	"testing"

	"github.com/sxllwx/vulcanus/pkg/command"
	"github.com/sxllwx/vulcanus/pkg/command/local"
	"github.com/sxllwx/vulcanus/pkg/command/remote"
)

var l = log.New(os.Stdout, "test ", log.Lshortfile|log.Ltime)

func newSSHHost() (command.Interface, error) {

	return remote.NewClient(&remote.Config{
		Remote:         "192.168.240.101:22",
		User:           "root",
		PrivateKeyFile: "/home/scott/.remote/id_rsa",
	}, l)

}

func newLocalHost() (command.Interface, error) {

	return local.NewClient(&local.Config{}, l), nil

}

func newManager() *Manager {

	h, err := newSSHHost()
	if err != nil {
		panic(err)
	}

	return NewManager(h)
}

func TestCreateChain(t *testing.T) {

	m := newManager()

	const (
		MainChainName = "DESKTOP-SERVICES"
		DesktopA      = "DESKTOP-A"
	)

	if err := m.CreateChainForTable(NATTableName, MainChainName); err != nil {
		t.Fatal(err)
	}

	if err := m.AppendChainToParentChain(NATTableName, PREROUTINGChainName, MainChainName, "desktop services portal"); err != nil {
		t.Fatal(err)
	}

	if err := m.CreateChainForTable(NATTableName, DesktopA); err != nil {
		t.Fatal(err)
	}

	if err := m.AppendChainToParentChain(NATTableName, MainChainName, DesktopA, "desktop-A-service"); err != nil {
		t.Fatal(err)
	}

	if err := m.AppendDNATRuleToChain(DesktopA, "tcp", "3000", "192.168.240.98:3000", "the desktop a policy"); err != nil {
		t.Fatal(err)
	}

	if err := m.DeleteChain(NATTableName, MainChainName, DesktopA, "desktop-A-service"); err != nil {
		t.Fatal(err)
	}

	if err := m.DeleteChain(NATTableName, PREROUTINGChainName, MainChainName, "desktop services portal"); err != nil {
		t.Fatal(err)
	}
}

func TestManager_CheckRule(t *testing.T) {

	//m := newManager()

	//if err := m.checkRule(NATTableName, PREROUTINGChainName, "-m", "comment","--comment", m.am.wrapComment("desktop services por//tal") ,"-j", "DESKTOP-SERVICES"); err != nil {
	//	t.Fatal(err)
	//}
}
