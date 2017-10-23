package bot

import (
	"encoding/json"
	"io/ioutil"

	"github.com/codegangsta/cli"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Token            string  `yaml:"token" json:"token"`
	Watch            bool    `yaml:"watch" json:"watch"`
	ListenAddr       string  `yaml:"listen_addr" json:"listen_addr"`
	MetricsPath      string  `yaml:"metrics_path" json:"metrics_path"`
	DatabasePath     string  `yaml:"database_path" json:"database_path"`
	DisableMessenger bool    `yaml:"disable_messenger" json:"disable_messenger"`
	Agents           []Agent `yaml:"agents" json:"agents"`
}

func NewConfig() *Config {
	return &Config{
		Watch:            false,
		ListenAddr:       ":9109",
		MetricsPath:      "/metrics",
		DatabasePath:     "./database.db",
		DisableMessenger: false,
	}
}

func (c *Config) LoadFromFile(file string) error {
	configBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(configBytes), &c); err != nil {
		return err
	}
	return nil
}

func (c *Config) LoadFromContext(con *cli.Context) {
	if len(con.String("token")) != 0 {
		c.Token = con.String("token")
	}
	if con.Bool("watch") == true {
		c.Watch = con.Bool("watch")
	}
	if len(con.String("listen-address")) != 0 {
		c.ListenAddr = con.String("listen-address")
	}
	if len(con.String("metrics-path")) != 0 {
		c.MetricsPath = con.String("metrics-path")
	}
	if len(con.String("database")) != 0 {
		c.DatabasePath = con.String("database")
	}
	if con.Bool("disable-messenger") == true {
		c.DisableMessenger = con.Bool("disable-messenger")
	}
}

func (c *Config) ToString() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(bytes)
}
