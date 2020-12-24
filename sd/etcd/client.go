// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	"google.golang.org/grpc"
)

const (
	beatDuration = 3 * 500 * time.Millisecond
)

const (
	dialKeepAlive = 3 * time.Second
	dialTimeout   = 3 * time.Second
)

type Client interface {
	Register(s Service) error
	Deregister(s Service) error

	GetEntries(prefix string) ([]string, error)
	WatchPrefix(prefix string, ch chan struct{})

	LeaseID() int64
}

type client struct {
	cli            *clientv3.Client
	ctx            context.Context
	kv             clientv3.KV
	leaser         clientv3.Lease
	leaseID        clientv3.LeaseID
	leaseKeepAlive <-chan *clientv3.LeaseKeepAliveResponse
	watcher        clientv3.Watcher
	watchCancel    context.CancelFunc
	watchCtx       context.Context
}

type Options struct {
	CACert        string
	CertFile      string
	DialKeepAlive time.Duration
	DialOptions   []grpc.DialOption
	DialTimeout   time.Duration
	KeyFile       string
	Password      string
	Username      string
}

func NewClient(ctx context.Context, endpoints []string, options Options) (Client, error) {
	var cfg *tls.Config
	var err error

	if options.DialKeepAlive == 0 {
		options.DialKeepAlive = dialKeepAlive
	}

	if options.DialTimeout == 0 {
		options.DialTimeout = dialTimeout
	}

	if options.CertFile != "" && options.KeyFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      options.CertFile,
			KeyFile:       options.KeyFile,
			TrustedCAFile: options.CACert,
		}

		if cfg, err = tlsInfo.ClientConfig(); err != nil {
			return nil, errors.Wrap(err, "failed to config")
		}
	}

	cli, err := clientv3.New(clientv3.Config{
		Context:           ctx,
		Endpoints:         endpoints,
		DialTimeout:       options.DialTimeout,
		DialKeepAliveTime: options.DialKeepAlive,
		DialOptions:       options.DialOptions,
		TLS:               cfg,
		Username:          options.Username,
		Password:          options.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to new")
	}

	return &client{
		cli: cli,
		ctx: ctx,
		kv:  clientv3.NewKV(cli),
	}, nil
}

func (c *client) Register(s Service) error {
	var err error

	if s.Key == "" {
		return errors.New("invalid key")
	}

	if s.Value == "" {
		return errors.New("invalid value")
	}

	if c.leaser != nil {
		if err := c.leaser.Close(); err != nil {
			return errors.Wrap(err, "failed to close")
		}
	}

	c.leaser = clientv3.NewLease(c.cli)

	if c.watcher != nil {
		if err := c.watcher.Close(); err != nil {
			return errors.Wrap(err, "failed to close")
		}
	}

	c.watcher = clientv3.NewWatcher(c.cli)

	if c.kv == nil {
		c.kv = clientv3.NewKV(c.cli)
	}

	resp, err := c.leaser.Grant(c.ctx, int64(beatDuration.Seconds()))
	if err != nil {
		return errors.Wrap(err, "failed to grant")
	}

	c.leaseID = resp.ID

	if _, err = c.kv.Put(c.ctx, s.Key, s.Value, clientv3.WithLease(c.leaseID)); err != nil {
		return errors.Wrap(err, "failed to put")
	}

	if c.leaseKeepAlive, err = c.leaser.KeepAlive(c.ctx, c.leaseID); err != nil {
		return errors.Wrap(err, "failed to keep alive")
	}

	go func() {
		for {
			select {
			case r := <-c.leaseKeepAlive:
				if r == nil {
					return
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c client) Deregister(s Service) error {
	defer c.close()

	if s.Key == "" {
		return errors.New("invalid key")
	}

	if _, err := c.cli.Delete(c.ctx, s.Key, clientv3.WithIgnoreLease()); err != nil {
		return errors.Wrap(err, "failed to delete")
	}

	return nil
}

func (c client) GetEntries(key string) ([]string, error) {
	resp, err := c.kv.Get(c.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get")
	}

	entries := make([]string, len(resp.Kvs))
	for i, item := range resp.Kvs {
		entries[i] = string(item.Value)
	}

	return entries, nil
}

func (c *client) WatchPrefix(prefix string, ch chan struct{}) {
	c.watcher = clientv3.NewWatcher(c.cli)
	c.watchCtx, c.watchCancel = context.WithCancel(c.ctx)

	wch := c.watcher.Watch(c.watchCtx, prefix, clientv3.WithPrefix(), clientv3.WithRev(0))
	ch <- struct{}{}
	for item := range wch {
		if item.Canceled {
			return
		}
		ch <- struct{}{}
	}
}

func (c client) LeaseID() int64 {
	return int64(c.leaseID)
}

func (c client) close() {
	if c.leaser != nil {
		if err := c.leaser.Close(); err != nil {
			// PASS
		}
	}

	if c.watcher != nil {
		if err := c.watcher.Close(); err != nil {
			// PASS
		}
	}

	if c.watchCancel != nil {
		c.watchCancel()
	}
}
