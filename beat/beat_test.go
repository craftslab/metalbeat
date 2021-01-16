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
	"context"
	"testing"

	"github.com/craftslab/metalbeat/etcd"
	"github.com/stretchr/testify/assert"
)

func TestBeat(t *testing.T) {
	endpoint := "127.0.0.1:2379"
	e := etcd.New(context.Background(), []string{endpoint}, etcd.DefaultConfig())

	c := DefaultConfig()
	c.Host = "127.0.0.1"
	assert.NotEqual(t, nil, c)

	b := New(c, e)
	assert.NotEqual(t, nil, b)
}
