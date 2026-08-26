package domain

var DefaultConfig = &Config{
	RefreshInterval: 10,
	Notifications: Notify{
		Before: 5,
		Active: true,
	},
}

type Config struct {
	RefreshInterval int        `yaml:"refresh_interval"`
	Notifications   Notify     `yaml:"notifications"`
	Calendars       []Calendar `yaml:"calendars"`
}

type Notify struct {
	Before int  `yaml:"before"`
	Active bool `yaml:"active"`
}

type Calendar struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Enabled  bool   `yaml:"enabled"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	LastSync int64  `yaml:"last_sync"`
	Error    string `yaml:"error"`
}
