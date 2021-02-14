package main

import (
	"flag"
	"fmt"
	"github.com/Dementir/decard/internal/logger"
	"github.com/Dementir/decard/internal/server"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func loadConfigFromYaml(path string) (*config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read config file: "+path)
	}
	defer f.Close()

	var conf config
	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse yaml-file")
	}

	return &conf, nil
}

func main() {
	lg := logger.NewLogger("DEBUG")
	defer lg.Sync()

	configPath := flag.String("c", "config.yaml", "set config path")
	flag.Parse()

	config, err := loadConfigFromYaml(*configPath)
	if err != nil {
		lg.Fatal("parse config error: " + err.Error())
		return
	}

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	lg.Info("run server")
	if err := server.InitServer(addr, lg); err != nil {
		os.Exit(1)
	}
}
