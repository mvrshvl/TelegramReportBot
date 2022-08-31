package config

type Config struct {
	Token    string                        `yaml:"token"`
	Login    string                        `yaml:"login"`
	Password string                        `yaml:"password"`
	Channels map[string]map[string]Channel `yaml:"channels"`
	Database Database                      `yaml:"database"`
}

type Channel struct {
	ID string `yaml:"id"`
}

type Database struct {
	Address  string `default:"vh354.timeweb.ru:3307"`
	User     string `default:"clubrums_bot"`
	Password string `default:"8Qq1SP4V"`
	Name     string `default:"clubrums_bot"`
	Driver   string `default:"mysql"`
}

func New() (*Config, error) {
	return &Config{
		Token:    "5440030440:AAGiPBn9JHb7VW60aRhf5S6HUqw6V8eEVgI",
		Login:    "clubrums_bot",
		Password: "8Qq1SP4V",
		Database: Database{
			Address:  "localhost",
			User:     "clubrums_bot",
			Password: "8Qq1SP4V",
			Name:     "clubrums_bot",
			Driver:   "mysql",
		},
		Channels: map[string]map[string]Channel{
			"volgograd": {
				"Охрана": Channel{ID: "-1001675961636"},
				"Бар":    Channel{ID: "-1001749986535"},
				"Арт":    Channel{ID: "-1001709159006"},
				"Админ":  Channel{ID: "-1001699370600"},
			},
			"moscow": {
				"Охрана": Channel{ID: "-1001697572469"},
				"Бар":    Channel{ID: "-1001605922666"},
				"Арт":    Channel{ID: "-1001770807023"},
				"Админ":  Channel{ID: "-1001579885451"},
			},
			"krasnodar": {
				"Охрана": Channel{ID: "-1001633711454"},
				"Бар":    Channel{ID: "-1001636988740"},
				"Арт":    Channel{ID: "-1001799331606"},
				"Админ":  Channel{ID: "-1001678631440"},
			},
		},
	}, nil
}
