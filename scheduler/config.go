package scheduler

// Config is struct for config.toml file
type Config struct {
	Port       string
	Bind       string
	AdminPort  string
	Suspend    int
	SocketPath string
	Slack      slackConfig
	Servers    []server
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
