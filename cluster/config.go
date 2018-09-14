package cluster

type Config struct {
	Master string
	//Port   string
}

func NewConfig(contents map[string]string) *Config {

	return &Config{
		Master: contents["master"],
		//Port:   contents["port"],
	}
}
