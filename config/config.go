package config

import "time"

//  "postgres://user:password@localhost:5432/mydb"
type Config struct {
	Driver          string
	User            string
	Password        string
	Host            string
	Port            string
	Database        string
	SemaphoreLimit  int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

var Cfg *Config

func SetConfig(c *Config) {
	Cfg = c
}

func GetConfig() *Config {
	if Cfg == nil {
		return &Config{
			Driver:          "postgres",
			User:            "postgres",
			Password:        "password",
			Host:            "localhost",
			Port:            "25432",
			Database:        "vulcan_test",
			SemaphoreLimit:  100,
			MaxOpenConns:    100,
			MaxIdleConns:    100,
			ConnMaxLifetime: 30 * time.Minute, // 30 min.
			ConnMaxIdleTime: 5 * time.Minute,  // 5 min.
		}
	} else {
		return Cfg
	}
}
