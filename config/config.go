package config

//  "postgres://user:password@localhost:5432/mydb"
type Config struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

var Cfg *Config

func SetConfig(c *Config) {
	Cfg = c
}

func GetConfig() *Config {
	if Cfg == nil {
		return &Config{
			Driver:   "postgres",
			User:     "amt",
			Password: "amt",
			Host:     "localhost",
			Port:     "25432",
			Database: "amt_consult",
		}
	} else {
		return Cfg
	}
}
