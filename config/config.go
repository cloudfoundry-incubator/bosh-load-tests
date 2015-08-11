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
