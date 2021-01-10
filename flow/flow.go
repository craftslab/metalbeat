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

package flow

import (
	"github.com/craftslab/metalbeat/beat"
	"github.com/craftslab/metalbeat/etcd"
)

type Flow interface {
	Run() error
}

type Config struct {
	Host string
}

type flow struct {
	beat beat.Beat
	etcd etcd.Etcd
	quit chan struct{}
	reg  string
	wch  string
}

func New(b beat.Beat, e etcd.Etcd, c *Config) Flow {
	return &flow{
		beat: b,
		etcd: e,
		reg:  "/metalflow/agent/" + c.Host + "/name",
		wch:  "/metalflow/worker/" + c.Host,
	}
}

func DefaultConfig() *Config {
	return &Config{
		Host: "127.0.0.1",
	}
}

func (f *flow) Run() error {
	ch := make(chan struct{})

	go func() {
		_ = f.etcd.Watch(f.wch, ch)
	}()

	defer func() {
		_ = f.etcd.Dewatch(f.wch)
	}()

	for {
		select {
		case <-ch:
			_, err := f.etcd.GetEntries(f.wch)
			if err != nil {
				continue
			}
		case <-f.quit:
			return nil
		}
	}
}
