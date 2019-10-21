package etcdv3

import (
	"context"
	"log"
	"path"
	"sync"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sxllwx/vulcanus/pkg/storage"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

const (
	defaultWatchBuf = 1
)

var (
	ErrNilETCDV3Client = errors.New("etcd raw client is nil") // full describe the ERR
	ErrKVPairNotFound  = errors.New("k/v pair not found")
)

type Options struct {
	name      string
	endpoints []string
	timeout   time.Duration
	heartbeat int // heartbeat second

	logger *log.Logger
}

type Option func(*Options)

func WithEndpoints(endpoints ...string) Option {
	return func(opt *Options) {
		opt.endpoints = endpoints
	}
}
func WithName(name string) Option {
	return func(opt *Options) {
		opt.name = name
	}
}
func WithTimeout(timeout time.Duration) Option {
	return func(opt *Options) {
		opt.timeout = timeout
	}
}

func WithHeartbeat(heartbeat int) Option {
	return func(opt *Options) {
		opt.heartbeat = heartbeat
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(opt *Options) {
		opt.logger = logger
	}
}

type Client struct {
	lock sync.RWMutex

	// these properties are only set once when they are started.
	name      string
	endpoints []string
	timeout   time.Duration
	heartbeat int

	logger *log.Logger

	ctx       context.Context    // if etcd server connection lose, the ctx.Done will be sent msg
	cancel    context.CancelFunc // cancel the ctx,  all watcher will stopped
	rawClient *clientv3.Client

	exit chan struct{}
	Wait sync.WaitGroup
}

func (c *Client) Reset() error {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.rawClient == nil {
		return ErrNilETCDV3Client
	}

	_, err := c.rawClient.Delete(c.ctx, "", clientv3.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(k string) error {

	err := c.delete(k)
	if err != nil {
		return errors.Annotatef(err, "delete k/v (key %s)", k)
	}

	return nil
}

func (c *Client) PUT(key string, data []byte) error {
	err := c.put(key, string(data))
	if err != nil {
		return errors.Annotatef(err, "put k/v (key: %s value %s)", key, data)
	}
	return nil
}

func (c *Client) GET(key string) ([]byte, error) {

	v, err := c.get(key)
	if err != nil {
		return nil, errors.Annotatef(err, "get key value (key %s)", key)
	}
	return []byte(v), nil
}

func (c *Client) WATCH(key string) (<-chan *storage.Message, error) {

	out := make(chan *storage.Message, defaultWatchBuf)

	wc, err := c.watch(key)
	if err != nil {
		return nil, errors.Annotatef(err, "watch prefix (key %s)", key)
	}

	go func() {

		for o := range wc {

			// got events
			for _, e := range o.Events {

				m := &storage.Message{
					Data: e.Kv.Value,
				}
				switch {

				case e.IsCreate():
					m.EventType = storage.Create
				case e.Type == mvccpb.DELETE:
					m.EventType = storage.Delete
				default:
					m.EventType = storage.Update
				}
				out <- m
			}
		}

		// the wc chan closed
		close(out)
	}()

	return out, nil
}

func (c *Client) List(prefix string) (map[string][]byte, error) {

	kList, vList, err := c.getChildren(prefix)
	if err != nil {
		return nil, errors.Annotatef(err, "get key children (key %s)", prefix)
	}

	out := map[string][]byte{}

	for i := 0; i < len(kList); i++ {
		out[kList[i]] = []byte(vList[i])
	}

	return out, nil
}

func (c *Client) Close() error {
	// stop the client
	if c.stop() {
		return nil
	}

	// wait client maintenance status stop
	c.Wait.Wait()

	c.lock.Lock()
	if c.rawClient != nil {
		c.clean()
	}
	c.lock.Unlock()
	c.logger.Printf("etcd client{name:%s, endpoints:%s} exit now.", c.name, c.endpoints)
	return nil
}

func (c *Client) Locker(k string) (storage.Locker, error) {

	session, err := concurrency.NewSession(c.rawClient)
	if err != nil {
		return nil, errors.Annotate(err, "new session")
	}
	m := concurrency.NewMutex(session, k)
	return m, nil
}

func NewEtcdV3Storage(opts ...Option) (storage.Interface, error) {

	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	ctx, cancel := context.WithCancel(context.Background())
	rawClient, err := clientv3.New(clientv3.Config{
		Context:     ctx,
		Endpoints:   o.endpoints,
		DialTimeout: o.timeout,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		return nil, errors.Annotate(err, "new raw client block connect to server")
	}

	c := &Client{

		name:      o.name,
		timeout:   o.timeout,
		endpoints: o.endpoints,
		heartbeat: o.heartbeat,

		logger: o.logger,

		ctx:       ctx,
		cancel:    cancel,
		rawClient: rawClient,

		exit: make(chan struct{}),
	}

	if err := c.maintenanceStatus(); err != nil {
		return nil, errors.Annotate(err, "client maintenance status")
	}
	return c, nil
}

// NOTICE: need to get the lock before calling this method
func (c *Client) clean() {

	// close raw client
	c.rawClient.Close()

	// cancel ctx for raw client
	c.cancel()

	// clean raw client
	c.rawClient = nil
}

func (c *Client) stop() bool {

	select {
	case <-c.exit:
		return true
	default:
		close(c.exit)
	}
	return false
}

func (c *Client) close() {

	// stop the client
	if c.stop() {
		return
	}

	// wait client maintenance status stop
	c.Wait.Wait()

	c.lock.Lock()
	if c.rawClient != nil {
		c.clean()
	}
	c.lock.Unlock()
	c.logger.Printf("etcd client{name:%s, endpoints:%s} exit now.", c.name, c.endpoints)
}

func (c *Client) maintenanceStatus() error {

	s, err := concurrency.NewSession(c.rawClient, concurrency.WithTTL(c.heartbeat))
	if err != nil {
		return errors.Annotate(err, "new session with server")
	}

	// must add wg before go maintenance status goroutine
	c.Wait.Add(1)
	go c.maintenanceStatusLoop(s)
	return nil
}

func (c *Client) maintenanceStatusLoop(s *concurrency.Session) {

	defer func() {
		c.Wait.Done()
		c.logger.Printf("etcd client {endpoints:%v, name:%s} maintenance goroutine game over.", c.endpoints, c.name)
	}()

	for {
		select {
		case <-c.Done():
			// Client be stopped, will clean the client hold resources
			return
		case <-s.Done():
			c.logger.Println("etcd server stopped")
			c.lock.Lock()
			// when etcd server stopped, cancel ctx, stop all watchers
			c.clean()
			// when connection lose, stop client, trigger reconnect to etcd
			c.stop()
			c.lock.Unlock()
			return
		}
	}
}

// if k not exist will put k/v in etcd
// if k is already exist in etcd, return nil
func (c *Client) put(k string, v string, opts ...clientv3.OpOption) error {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.rawClient == nil {
		return ErrNilETCDV3Client
	}

	_, err := c.rawClient.Txn(c.ctx).
		If(clientv3.Compare(clientv3.Version(k), "<", 1)).
		Then(clientv3.OpPut(k, v, opts...)).
		Commit()
	if err != nil {
		return err

	}
	return nil
}

func (c *Client) delete(k string) error {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.rawClient == nil {
		return ErrNilETCDV3Client
	}

	_, err := c.rawClient.Delete(c.ctx, k)
	if err != nil {
		return err

	}
	return nil
}

func (c *Client) get(k string) (string, error) {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.rawClient == nil {
		return "", ErrNilETCDV3Client
	}

	resp, err := c.rawClient.Get(c.ctx, k)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", ErrKVPairNotFound
	}

	return string(resp.Kvs[0].Value), nil
}

func (c *Client) getChildren(k string) ([]string, []string, error) {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.rawClient == nil {
		return nil, nil, ErrNilETCDV3Client
	}

	resp, err := c.rawClient.Get(c.ctx, k, clientv3.WithPrefix())
	if err != nil {
		return nil, nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, nil, ErrKVPairNotFound
	}

	var (
		kList []string
		vList []string
	)

	for _, kv := range resp.Kvs {
		kList = append(kList, string(kv.Key))
		vList = append(vList, string(kv.Value))
	}

	return kList, vList, nil
}

func (c *Client) watchWithPrefix(prefix string) (clientv3.WatchChan, error) {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.rawClient == nil {
		return nil, ErrNilETCDV3Client
	}

	return c.rawClient.Watch(c.ctx, prefix, clientv3.WithPrefix()), nil
}

func (c *Client) watch(k string) (clientv3.WatchChan, error) {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.rawClient == nil {
		return nil, ErrNilETCDV3Client
	}

	return c.rawClient.Watch(c.ctx, k), nil
}

func (c *Client) keepAliveKV(k string, v string) error {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.rawClient == nil {
		return ErrNilETCDV3Client
	}

	lease, err := c.rawClient.Grant(c.ctx, int64(time.Second.Seconds()))
	if err != nil {
		return errors.Annotate(err, "grant lease")
	}

	keepAlive, err := c.rawClient.KeepAlive(c.ctx, lease.ID)
	if err != nil || keepAlive == nil {
		c.rawClient.Revoke(c.ctx, lease.ID)
		return errors.Annotate(err, "keep alive lease")
	}

	_, err = c.rawClient.Put(c.ctx, k, v, clientv3.WithLease(lease.ID))
	if err != nil {
		return errors.Annotate(err, "put k/v with lease")
	}
	return nil
}

func (c *Client) Done() <-chan struct{} {
	return c.exit
}

func (c *Client) Valid() bool {
	select {
	case <-c.exit:
		return false
	default:
	}

	c.lock.RLock()
	if c.rawClient == nil {
		c.lock.RUnlock()
		return false
	}
	c.lock.RUnlock()
	return true
}

func (c *Client) RegisterTemp(basePath string, node string) (string, error) {

	completeKey := path.Join(basePath, node)

	err := c.keepAliveKV(completeKey, "")
	if err != nil {
		return "", errors.Annotatef(err, "keepalive kv (key %s)", completeKey)
	}

	return completeKey, nil
}
