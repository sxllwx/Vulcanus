package remote

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"github.com/sxllwx/vulcanus/pkg/command"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	cfg     *Config
	session *ssh.Session
	client  *ssh.Client
	logger  *log.Logger
}

type Config struct {
	Remote         string
	User           string
	PrivateKeyFile string
}

func NewClient(cfg *Config) (command.Interface, error) {

	key, err := ioutil.ReadFile(cfg.PrivateKeyFile)
	if err != nil {
		return nil, errors.WithMessagef(err, "read private key file %s", cfg.PrivateKeyFile)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse the private key")
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
		return nil, errors.WithMessage(err, "dial")
	}

	return &Client{
		cfg:    cfg,
		client: clt,
	}, nil

}

func (c *Client) Close() error {

	c.client.Close()
	return nil
}

func (c *Client) Exec(cmd string, args []string, in io.Reader, out, err io.WriteCloser) error {

	buff := &bytes.Buffer{}
	buff.WriteString(cmd)
	for _, a := range args {
		buff.WriteString(" ")
		buff.WriteString(a)
	}

	s, e := c.client.NewSession()
	if e != nil {
		return errors.WithMessage(e, "new session")
	}
	defer s.Close()

	s.Stderr = err
	s.Stdin = in
	s.Stdout = out

	e = s.Run(buff.String())
	if e != nil {
		return errors.WithMessage(e, "run cmd")
	}

	return nil
}
