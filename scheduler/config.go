package scheduler

// Config is struct for config.toml file
type Config struct {
	Port       string
	Bind       string
	AdminPort  string
	Suspend    int
	SocketPath string
	Servers    []server
}

type server struct {
	Hostname string
	IP       string
	Port     string
	Weight   int
	Backup   bool
	Token    string
}
