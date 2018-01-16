package scheduler

// Config is struct for config.toml file
type Config struct {
	Port      int
	Bind      string
	AdminPort int
	Suspend   int
	Servers   []server
}

type server struct {
	IP     string
	Port   string
	Weight int
	Backup bool
	Token  string
}
