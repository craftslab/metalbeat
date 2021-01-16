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
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/craftslab/metalbeat/etcd"
	"github.com/craftslab/metalbeat/runtime"
)

const (
	agent  = "/metalflow/agent/"
	worker = "/metalflow/worker/"
)

const (
	sep = " "
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
		reg:  agent + c.Host + "/register",
		wch:  worker + c.Host,
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
	ch := make(chan map[string]string)

	go func() {
		_ = b.etcd.Watch(b.wch, ch)
	}()

	for {
		select {
		case ev := <-ch:
			if key, ok := ev[etcd.EventPut]; ok {
				entries, err := b.etcd.GetEntries(key)
				if err != nil {
					return errors.Wrap(err, "failed to get")
				}
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
	for _, item := range event {
		buf := []interface{}{
			item,
		}
		if _, err := runtime.Run(b.routine, buf); err != nil {
			return errors.Wrap(err, "failed to run")
		}
	}

	return nil
}

func (b *beat) routine(cmd interface{}) interface{} {
	slice := strings.Split(cmd.(string), sep)

	c := exec.Command(slice[0], slice[1:]...) // nolint:gosec
	_, err := c.CombinedOutput()

	return err
}
