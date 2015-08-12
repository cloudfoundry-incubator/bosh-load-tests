package config

import (
	"encoding/json"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Config struct {
	Environment              string `json:"environment"`
	DirectorMigrationCommand string `json:"director_migration_cmd"`
	DirectorStartCommand     string `json:"director_start_cmd"`
	WorkerStartCommand       string `json:"worker_start_cmd"`
	NginxStartCommand        string `json:"nginx_start_cmd"`
	NatsStartCommand         string `json:"nats_start_cmd"`
	CliCmd                   string `json:"cli_cmd"`
	NumberOfDeployments      int    `json:"number_of_deployments"`
	fs                       boshsys.FileSystem
}

func NewConfig(fs boshsys.FileSystem) *Config {
	return &Config{
		fs: fs,
	}
}

func (c *Config) Load(configPath string) error {
	contents, err := c.fs.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(contents), &c)
	if err != nil {
		return err
	}

	return nil
}
