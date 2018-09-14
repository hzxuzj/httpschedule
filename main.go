package main

import (
	"fmt"
	"log"

	"httpschedule2/builtins"
	"httpschedule2/daemon"
	"httpschedule2/engine"
)

func main() {
	eng := engine.NewEngine()

	builtins.Register(eng)

	d, err := daemon.NewDaemon(eng)

	if err != nil {
		fmt.Printf("err : %v\n", err)
		fmt.Errorf("%s", err)
		return
	}
	d.Install(eng)

	flHost := []string{d.Conf.ProtoAddr}
	job := eng.Job("serverapi", flHost...)
	job.SetEnvBool("EnableCors", true)
	if err := job.Run(); err != nil {
		log.Fatal(err)
	}
}
