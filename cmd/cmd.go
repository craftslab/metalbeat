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
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"

	"github.com/craftslab/metalbeat/beat"
	"github.com/craftslab/metalbeat/config"
	"github.com/craftslab/metalbeat/etcd"
	"github.com/craftslab/metalbeat/service"
)

var (
	app        = kingpin.New("metalbeat", "Metal Beat").Version(config.Version + "-build-" + config.Build)
	configFile = app.Flag("config-file", "Config file (.yml)").Required().String()
	hostAddr   = app.Flag("host-addr", "Host address").Required().String()
)

func Run() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	c, err := initConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	e, err := initEtcd(c)
	if err != nil {
		log.Fatalf("failed to init etcd: %v", err)
	}

	b, err := initBeat(c)
	if err != nil {
		log.Fatalf("failed to init beat: %v", err)
	}

	log.Println("beat running")

	if err := runService(c, b, e); err != nil {
		log.Fatalf("failed to run beat: %v", err)
	}

	log.Println("beat exiting")
}

func initConfig(name string) (*config.Config, error) {
	c := config.New()
	if c == nil {
		return &config.Config{}, errors.New("failed to new")
	}

	fi, err := os.Open(name)
	if err != nil {
		return c, errors.Wrap(err, "failed to open")
	}

	defer func() {
		_ = fi.Close()
	}()

	buf, err := ioutil.ReadAll(fi)
	if err != nil {
		return c, errors.Wrap(err, "failed to readall")
	}

	if err := yaml.Unmarshal(buf, c); err != nil {
		return c, errors.Wrap(err, "failed to unmarshal")
	}

	return c, nil
}

func initEtcd(cfg *config.Config) (etcd.Etcd, error) {
	endpoint := cfg.Spec.Sd.Host + ":" + cfg.Spec.Sd.Port
	return etcd.New(context.Background(), []string{endpoint}, etcd.DefaultConfig()), nil
}

func initBeat(_ *config.Config) (beat.Beat, error) {
	return beat.New(beat.DefaultConfig()), nil
}

func runService(_ *config.Config, b beat.Beat, e etcd.Etcd) error {
	cfg := service.DefaultConfig()
	if cfg == nil {
		return errors.New("failed to config")
	}

	cfg.Host = *hostAddr

	srv := service.New(b, e, cfg)
	if srv == nil {
		return errors.New("failed to new")
	}

	return srv.Run()
}
