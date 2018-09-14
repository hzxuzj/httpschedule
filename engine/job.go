package engine

import (
	"fmt"
	"io"
	"time"
)

type Status int

const StatusOK Status = 0
const StatusErr Status = -1

type Job struct {
	Name    string
	Eng     *Engine
	Args    []string
	handler Handler
	end     time.Time
	Stdin   *Input
	Stdout  *Output
	Stderr  *Output
	env     *Env
	status  Status
}

func (job *Job) Run() error {

	if job.Name != "serverapi" {
		job.Eng.l.Lock()
		job.Eng.tasks.Add(1)
		job.Eng.l.Unlock()
		defer job.Eng.tasks.Done()
	}

	if !job.end.IsZero() {
		return fmt.Errorf("%s: job has already completed", job.Name)
	}

	if job.handler == nil {
		return fmt.Errorf("%s: command not found", job.Name)
	} else {
		job.status = job.handler(job)
		job.end = time.Now()
	}

	if err := job.Stdout.Close(); err != nil {
		return err
	}

	if err := job.Stderr.Close(); err != nil {
		return err
	}

	if err := job.Stdin.Close(); err != nil {
		return err
	}

	if job.status != StatusOK {
		return fmt.Errorf("job failed")
	}

	return nil
}

func (job *Job) GetEnvJson(key string, iface interface{}) error {
	return job.env.GetJson(key, iface)
}

func (job *Job) SetEnvJson(key string, value interface{}) error {
	return job.env.SetJson(key, value)
}

func (job *Job) SetEnv(key, value string) {
	job.env.Set(key, value)
}

func (job *Job) GetEnv(key string) string {
	return job.env.Get(key)
}

func (job *Job) SetEnvBool(key string, value bool) {
	job.env.SetBool(key, value)
}

func (job *Job) GetEnvBool(key string) bool {
	return job.env.GetBool(key)
}

func (job *Job) SetEnvInt(key string, value int) {
	job.env.SetInt(key, value)
}

func (job *Job) GetEnvInt(key string) int {
	return job.env.GetInt(key)
}

func (job *Job) SetEnvInt64(key string, value int64) {
	job.env.SetInt64(key, value)
}

func (job *Job) GetEnvInt64(key string) int64 {
	return job.env.GetInt64(key)
}

func (job *Job) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(job.Stdout, format, args)
}

func (job *Job) Errorf(format string, args ...interface{}) Status {

	if format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	fmt.Fprintf(job.Stderr, format, args)

	return StatusErr
}
func (job *Job) DecodeEnv(src io.Reader) error {
	return job.env.Decode(src)
}

func (job *Job) WriteError(err error) error {
	m := newMessage(err.Error())
	if _, er := job.Stdout.Write(m.toJsonBytes()); er != nil {
		return er
	}
	return nil
}

func (job *Job) WriteOK() error {
	m := newMessage("ok")
	if _, err := job.Stdout.Write(m.toJsonBytes()); err != nil {
		return err
	}

	return nil
}
