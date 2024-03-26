package config

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	config     Config
	once       sync.Once
	configFile string = "config.yaml"
)

type Config struct {
	Port         int           `yaml:"port"`
	LogLevel     string        `yaml:"log_level"`
	LogTimestamp bool          `yaml:"log_timestamp"`
	Interval     time.Duration `yaml:"interval"`
	Monitors     []*Monitor
}

type Monitor struct {
	Name     string
	Items    []*Item
	Webhooks []*Webhook
}

type Item struct {
	Name      string
	URL       string
	Regex     string
	Value     float64
	LastCheck time.Time `yaml:"last_check"`
}

type Webhook struct {
	URL      string
	Template string
}

func Get() *Config {
	once.Do(func() {
		// initialize
		config = Config{
			Port:         8080,
			LogLevel:     "info",
			LogTimestamp: false,
			Monitors:     make([]*Monitor, 0),
		}

		// load config file
		if _, err := os.Stat(configFile); err == nil {
			data, err := ioutil.ReadFile(configFile)
			if err != nil {
				log.Println("could not load", configFile)
				log.Fatalln(err.Error())
			}
			if err := yaml.Unmarshal(data, &config); err != nil {
				log.Println("could not parse", configFile)
				log.Fatalln(err.Error())
			}
		} else {
			log.Println("could not find", configFile)
			log.Fatalln(err.Error())
		}
	})
	return &config
}
