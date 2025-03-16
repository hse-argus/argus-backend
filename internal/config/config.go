package config

type Config struct {
	WebAddr string

	PostgresDb       string
	PostgresUser     string
	PostgresPassword string
	PostgresPort     int
}

func NewConfig() *Config {
	return &Config{
		WebAddr:          ":3000",
		PostgresDb:       "argus",
		PostgresUser:     "username",
		PostgresPassword: "1234",
		PostgresPort:     5432,
	}
}
