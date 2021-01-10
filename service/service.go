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

package service

import (
	"github.com/craftslab/metalbeat/beat"
	"github.com/craftslab/metalbeat/etcd"
)

type Service interface {
	Run() error
}

type Config struct {
	Host string
}

type service struct {
	beat beat.Beat
	etcd etcd.Etcd
	quit chan struct{}
	reg  string
	wch  string
}

func New(b beat.Beat, e etcd.Etcd, c *Config) Service {
	return &service{
		beat: b,
		etcd: e,
		reg:  "/metalbeat/" + c.Host,
		wch:  "/metalflow/" + c.Host,
	}
}

func DefaultConfig() *Config {
	return &Config{
		Host: "127.0.0.1",
	}
}

func (b *service) Run() error {
	ch := make(chan struct{})

	go func() {
		_ = b.etcd.Watch(b.wch, ch)
	}()

	defer func() {
		_ = b.etcd.Dewatch(b.wch)
	}()

	for {
		select {
		case <-ch:
			_, err := b.etcd.GetEntries(b.wch)
			if err != nil {
				continue
			}
		case <-b.quit:
			return nil
		}
	}
}
