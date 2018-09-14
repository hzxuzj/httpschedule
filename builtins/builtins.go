package builtins

import (
	"httpschedule2/api"
	"httpschedule2/engine"
)

func Register(eng *engine.Engine) error {
	if err := eng.Register("serverapi", api.ServerApi); err != nil {
		return err
	}

	return nil
}
