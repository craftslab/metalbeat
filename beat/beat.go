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

package beat

import (
	"github.com/craftslab/metalbeat/etcd"
	"github.com/pkg/errors"
)

type Beat interface {
	Run() error
}

type Config struct {
	Host string
}

type beat struct {
	etcd etcd.Etcd
	quit chan struct{}
	reg  string
	wch  string
}

func New(c *Config, e etcd.Etcd) Beat {
	return &beat{
		etcd: e,
		reg:  "/metalflow/agent/" + c.Host + "/register",
		wch:  "/metalflow/worker/" + c.Host,
	}
}

func DefaultConfig() *Config {
	return &Config{
		Host: "127.0.0.1",
	}
}

func (b *beat) Run() error {
	if err := b.register(); err != nil {
		return errors.Wrap(err, "failed to register")
	}

	if err := b.watch(); err != nil {
		return errors.Wrap(err, "failed to watch")
	}

	return nil
}

func (b *beat) register() error {
	return b.etcd.Register(b.reg, "metalbeat", etcd.TTLDuration)
}

func (b *beat) watch() error {
	ch := make(chan struct{})

	go func() {
		_ = b.etcd.Watch(b.wch, ch)
	}()

	for {
		select {
		case <-ch:
			if entries, err := b.etcd.GetEntries(b.wch); err == nil {
				if err := b.dispatch(entries); err != nil {
					return errors.Wrap(err, "failed to dispatch")
				}
			}
		case <-b.quit:
			return nil
		}
	}
}

func (b *beat) dispatch(event []string) error {
	// TODO
	return nil
}
