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
	dialKeepAlive = 3 * time.Second
	dialTimeout   = 3 * time.Second
)

const (
	ttlDuration = 3 * 500 * time.Millisecond
)

type Etcd interface {
	Register(key, val string, ttl time.Duration) error
	Deregister(key string) error

	GetEntries(prefix string) ([]string, error)
	WatchPrefix(prefix string, ch chan struct{})

	LeaseID() int64
}

type Option struct {
	CACert        string
	CertFile      string
	DialKeepAlive time.Duration
	DialOptions   []grpc.DialOption
	DialTimeout   time.Duration
	KeyFile       string
	Password      string
	Username      string
}

type etcd struct {
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

func New(ctx context.Context, endpoints []string, option Option) Etcd {
	var cfg *tls.Config
	var err error

	if option.DialKeepAlive == 0 {
		option.DialKeepAlive = dialKeepAlive
	}

	if option.DialTimeout == 0 {
		option.DialTimeout = dialTimeout
	}

	if option.CertFile != "" && option.KeyFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      option.CertFile,
			KeyFile:       option.KeyFile,
			TrustedCAFile: option.CACert,
		}

		if cfg, err = tlsInfo.ClientConfig(); err != nil {
			return nil
		}
	}

	cli, err := clientv3.New(clientv3.Config{
		Context:           ctx,
		Endpoints:         endpoints,
		DialTimeout:       option.DialTimeout,
		DialKeepAliveTime: option.DialKeepAlive,
		DialOptions:       option.DialOptions,
		TLS:               cfg,
		Username:          option.Username,
		Password:          option.Password,
	})
	if err != nil {
		return nil
	}

	return &etcd{
		cli: cli,
		ctx: ctx,
		kv:  clientv3.NewKV(cli),
	}
}

func (e *etcd) Register(key, val string, ttl time.Duration) error {
	var err error

	if key == "" {
		return errors.New("invalid key")
	}

	if val == "" {
		return errors.New("invalid value")
	}

	if ttl == 0 {
		ttl = ttlDuration
	}

	if e.leaser != nil {
		if err := e.leaser.Close(); err != nil {
			return errors.Wrap(err, "failed to close")
		}
	}

	e.leaser = clientv3.NewLease(e.cli)

	if e.watcher != nil {
		if err := e.watcher.Close(); err != nil {
			return errors.Wrap(err, "failed to close")
		}
	}

	e.watcher = clientv3.NewWatcher(e.cli)

	if e.kv == nil {
		e.kv = clientv3.NewKV(e.cli)
	}

	resp, err := e.leaser.Grant(e.ctx, int64(ttl.Seconds()))
	if err != nil {
		return errors.Wrap(err, "failed to grant")
	}

	e.leaseID = resp.ID

	if _, err = e.kv.Put(e.ctx, key, val, clientv3.WithLease(e.leaseID)); err != nil {
		return errors.Wrap(err, "failed to put")
	}

	if e.leaseKeepAlive, err = e.leaser.KeepAlive(e.ctx, e.leaseID); err != nil {
		return errors.Wrap(err, "failed to keep alive")
	}

	go func() {
		for {
			select {
			case r := <-e.leaseKeepAlive:
				if r == nil {
					return
				}
			case <-e.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (e etcd) Deregister(key string) error {
	defer e.close()

	if key == "" {
		return errors.New("invalid key")
	}

	if _, err := e.cli.Delete(e.ctx, key, clientv3.WithIgnoreLease()); err != nil {
		return errors.Wrap(err, "failed to delete")
	}

	return nil
}

func (e etcd) GetEntries(key string) ([]string, error) {
	resp, err := e.kv.Get(e.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get")
	}

	entries := make([]string, len(resp.Kvs))
	for i, item := range resp.Kvs {
		entries[i] = string(item.Value)
	}

	return entries, nil
}

func (e *etcd) WatchPrefix(prefix string, ch chan struct{}) {
	e.watcher = clientv3.NewWatcher(e.cli)
	e.watchCtx, e.watchCancel = context.WithCancel(e.ctx)

	wch := e.watcher.Watch(e.watchCtx, prefix, clientv3.WithPrefix(), clientv3.WithRev(0))
	ch <- struct{}{}
	for item := range wch {
		if item.Canceled {
			return
		}
		ch <- struct{}{}
	}
}

func (e etcd) LeaseID() int64 {
	return int64(e.leaseID)
}

func (e etcd) close() {
	if e.leaser != nil {
		if err := e.leaser.Close(); err != nil {
			// PASS
		}
	}

	if e.watcher != nil {
		if err := e.watcher.Close(); err != nil {
			// PASS
		}
	}

	if e.watchCancel != nil {
		e.watchCancel()
	}
}
