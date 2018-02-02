package scheduler

// Config is struct for config.toml file
type Config struct {
	Port       string
	Bind       string
	AdminPort  string
	Suspend    int
	Ticker     int
	SocketPath string
	Slack      slackConfig
	Telegram   telegramConfig
	Servers    []server
}

type telegramConfig struct {
	Use   bool
	Token string
}

type slackConfig struct {
	Use     bool
	Token   string
	Channel string
}

type server struct {
	Hostname string
	IP       string
	Port     string
	Weight   int
	Backup   bool
	Token    string
}
