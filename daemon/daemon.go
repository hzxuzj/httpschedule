package daemon

import (
	"httpschedule2/cluster"
	"httpschedule2/engine"
)

type Daemon struct {
	Eng           *engine.Engine
	ApiVersion    string
	ClusterConfig *cluster.Config
	Conf          *Config
}

func NewDaemon(eng *engine.Engine) (*Daemon, error) {
	conf := &Config{}

	contents := make(map[string]string)
	contents["protoAddr"] = "tcp://0.0.0.0:9090"

	conf.ProtoAddr = contents["protoAddr"]

	clusterConfig := cluster.NewConfig(contents)
	daemon := &Daemon{
		Eng:           eng,
		ApiVersion:    "/api/v1",
		ClusterConfig: clusterConfig,
		Conf:          conf,
	}

	return daemon, nil
}

func (d *Daemon) Install(eng *engine.Engine) error {
	for name, method := range map[string]engine.Handler{
		"selectnode": d.selectnode,
	} {
		if err := eng.Register(name, method); err != nil {
			return err
		}
	}
	return nil

}
