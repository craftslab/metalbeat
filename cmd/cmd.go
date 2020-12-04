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

package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"

	"github.com/craftslab/metalbeat/config"
)

var (
	app        = kingpin.New("metalbeat", "Metal Beat").Version(config.Version + "-build-" + config.Build)
	configFile = app.Flag("config-file", "Config file (.yml)").Required().String()
)

var (
	cfg = config.Config{}
)

func Run() {
	var err error

	kingpin.MustParse(app.Parse(os.Args[1:]))

	cfg, err = parseConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to parse: %v", err)
	}

	log.Println("beat running")
	log.Println("beat exiting")
}

func parseConfig(name string) (config.Config, error) {
	var c config.Config

	fi, err := os.Open(name)
	if err != nil {
		return c, errors.Wrap(err, "failed to open")
	}

	defer func() {
		_ = fi.Close()
	}()

	buf, _ := ioutil.ReadAll(fi)
	if err := yaml.Unmarshal(buf, &c); err != nil {
		return c, errors.Wrap(err, "failed to unmarshal")
	}

	return c, nil
}
