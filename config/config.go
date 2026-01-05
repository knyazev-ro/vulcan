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

func GetConfig() *Config {
	return &Config{
		Driver:   "postgres",
		User:     "postgres",
		Password: "123",
		Host:     "localhost",
		Port:     "5432",
		Database: "vulcan_test",
	}
}
