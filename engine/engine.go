package engine

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Handler func(*Job) Status

type Engine struct {
	handlers map[string]Handler
	tasks    sync.WaitGroup
	l        sync.RWMutex
	shutdown []func()
	Stdout   io.Writer
	Stderr   io.Writer
	Stdin    io.Reader
}

func (engine *Engine) Register(name string, method Handler) error {
	_, exists := engine.handlers[name]

	if exists {
		return fmt.Errorf("%s is already exists!")
	}

	engine.handlers[name] = method
	return nil
}

func NewEngine() *Engine {
	eng := &Engine{
		handlers: make(map[string]Handler),
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Stdin:    os.Stdin,
	}

	return eng
}

func (eng *Engine) Job(name string, args ...string) *Job {
	job := &Job{
		Eng:    eng,
		Name:   name,
		Args:   args,
		Stdin:  NewInput(),
		Stdout: NewOutput(),
		Stderr: NewOutput(),
		env:    &Env{},
	}

	if handler, exists := eng.handlers[name]; exists {
		job.handler = handler
	}
	return job
}
