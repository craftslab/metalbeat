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
	"github.com/craftslab/metalbeat/context"
)

type Flow struct {
	Config Config
}

type Config struct {
}

func New(cfg Config) *Flow {
	return &Flow{
		Config: cfg,
	}
}

func DefaultConfig() Config {
	return Config{}
}

func (f Flow) Use(ctx context.Context) error {
	return nil
}

func (f Flow) Run() error {
	return nil
}
