package domain

var DefaultConfig = &Config{
	TemporaryFilesPath: "tmp/",
	RefreshInterval:    10,
	Notifications: Notify{
		Before: 5,
		Active: true,
	},
}

type Config struct {
	RefreshInterval    int        `yaml:"refresh_interval"`
	Notifications      Notify     `yaml:"notifications"`
	Calendars          []Calendar `yaml:"calendars"`
	TemporaryFilesPath string     `yaml:"temporary_files_path"`
}

type Notify struct {
	Before int  `yaml:"before"`
	Active bool `yaml:"active"`
}

type Calendar struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Enabled  bool   `yaml:"enabled"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
