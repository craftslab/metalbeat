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
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestEtcd(t *testing.T) {
	config := Config{
		CACert:        "",
		CertFile:      "",
		DialKeepAlive: DialKeepAlive,
		DialOptions:   []grpc.DialOption{grpc.WithBlock()},
		DialTimeout:   DialTimeout,
		KeyFile:       "",
		Password:      "",
		Username:      "",
	}

	endpoint := "127.0.0.1:2379"

	e := New(context.Background(), []string{endpoint}, &config)
	assert.NotEqual(t, nil, e)

	err := e.Register("", "", TTLDuration)
	assert.NotEqual(t, nil, err)

	err = e.Watch("", nil)
	assert.NotEqual(t, nil, err)

	err = e.Dewatch("")
	assert.NotEqual(t, nil, err)

	_, err = e.GetEntries("")
	assert.NotEqual(t, nil, err)

	err = e.Deregister("")
	assert.NotEqual(t, nil, err)

	prefix := "/metalflow/agent"
	key := prefix + "/127.0.0.1/name"
	val := "metalbeat"

	err = e.Register(key, val, TTLDuration)
	assert.Equal(t, nil, err)

	entries, err := e.GetEntries(key)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(entries))

	id := e.LeaseID()
	assert.NotEqual(t, 0, id)

	err = e.Deregister(key)
	assert.Equal(t, nil, err)
}
