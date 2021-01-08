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
	endpoint := "127.0.0.1:2379"
	config := Config{
		CACert:        "",
		CertFile:      "",
		DialKeepAlive: dialKeepAlive,
		DialOptions:   []grpc.DialOption{grpc.WithBlock()},
		DialTimeout:   dialTimeout,
		KeyFile:       "",
		Password:      "",
		Username:      "",
	}

	e := New(context.Background(), []string{endpoint}, &config)
	assert.NotEqual(t, nil, e)

	key := "/metalflow/127.0.0.1"
	val := "metalbeat"

	err := e.Register(key, val, ttlDuration)
	assert.Equal(t, nil, err)

	_, err = e.GetEntries(key)
	assert.Equal(t, nil, err)

	id := e.LeaseID()
	assert.NotEqual(t, 0, id)

	err = e.Deregister(key)
	assert.Equal(t, nil, err)
}
