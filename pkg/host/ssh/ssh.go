package ssh

import (
	"bytes"
	"io/ioutil"

	"github.com/juju/errors"
	"github.com/sxllwx/vulcanus/pkg/host"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	cfg     *Config
	session *ssh.Session
}

type Config struct {
	Remote         string
	User           string
	PrivateKeyFile string
}

func NewClient(cfg *Config) (host.Interface, error) {

	key, err := ioutil.ReadFile(cfg.PrivateKeyFile)
	if err != nil {
		return nil, errors.Annotatef(err, "read private key file %s", cfg.PrivateKeyFile)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, errors.Annotate(err, "parse the private key")
	}

	clt, err := ssh.Dial("tcp", cfg.Remote, &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, errors.Annotate(err, "dial")
	}

	s, err := clt.NewSession()
	if err != nil {
		return nil, errors.Annotate(err, "new session")
	}

	return &Client{
		cfg:     cfg,
		session: s,
	}, nil

}

func (c *Client) Execute(rootCommand string, args ...string) ([]byte, error) {

	buff := &bytes.Buffer{}
	buff.WriteString(rootCommand)
	for _, a := range args {
		buff.WriteString(" ")
		buff.WriteString(a)
	}

	out, err := c.session.CombinedOutput(buff.String())
	if err != nil {
		return nil, errors.Annotate(err, "run cmd")
	}

	return out, nil
}
