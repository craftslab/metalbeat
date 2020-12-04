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

package config

type Config struct {
	MetaData MetaData `yaml:"metadata"`
	Spec     []Spec     `yaml:"spec"`
}

type MetaData struct {
	Name string `yaml:"name"`
}

type Spec struct {
	Name string                 `yaml:"name"`
	Node map[string]interface{} `yaml:"node"`
	Role string                 `yaml:"role"`
}

type Bare struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Container struct {
	Command []string `yaml:"command"`
	Expose  []int    `yaml:"expose"`
	Image   string   `yaml:"image"`
}

const (
	NodeBare      = "bare"
	NodeContainer = "container"
)

var (
	Build   string
	Version string
)
