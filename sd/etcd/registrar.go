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
	"sync"

	"github.com/pkg/errors"
)

type Registrar struct {
	client  Client
	mutex   sync.Mutex
	quit    chan struct{}
	service Service
}

type Service struct {
	Key   string
	Value string
}

func NewRegistrar(client Client, service Service) *Registrar {
	return &Registrar{
		client:  client,
		service: service,
	}
}

func (r *Registrar) Register() error {
	if err := r.client.Register(r.service); err != nil {
		return errors.Wrap(err, "failed to register")
	}

	return nil
}

func (r *Registrar) Deregister() error {
	if err := r.client.Deregister(r.service); err != nil {
		// PASS
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.quit != nil {
		close(r.quit)
		r.quit = nil
	}

	return nil
}
