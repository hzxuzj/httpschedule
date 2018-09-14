package engine

import (
	"fmt"
	"io"
	"sync"
)

type Output struct {
	sync.Mutex
	dests []io.Writer
	//task  sync.WaitGroup
	used bool
}

func NewOutput() *Output {
	return &Output{}
}

func (o *Output) Close() error {
	o.Lock()
	defer o.Unlock()

	var firstErr error
	for _, dst := range o.dests {
		if closer, ok := dst.(io.Closer); ok {
			fmt.Println("close....")
			err := closer.Close()
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	//o.task.Wait()
	return firstErr
}

func (o *Output) Set(dst io.Writer) {
	o.Close()
	o.Lock()
	defer o.Unlock()
	o.dests = []io.Writer{dst}
}

func (o *Output) Add(dst io.Writer) {
	o.Lock()
	defer o.Unlock()
	o.dests = append(o.dests, dst)
}

func (o *Output) Write(p []byte) (n int, err error) {
	o.Lock()
	defer o.Unlock()
	o.used = true
	var firstErr error
	for _, dst := range o.dests {
		_, err := dst.Write(p)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return len(p), firstErr
}

func (o *Output) AddPipe() (io.Reader, error) {
	r, w := io.Pipe()
	o.Add(w)
	return r, nil
}

type Input struct {
	src io.Reader
	sync.Mutex
}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) Read(p []byte) (n int, err error) {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()
	if i.src == nil {
		return 0, io.EOF
	}
	return i.src.Read(p)
}

func (i *Input) Close() error {
	if i.src != nil {
		if closer, ok := i.src.(io.Closer); ok {
			return closer.Close()
		}
	}
	return nil
}

func (i *Input) Add(src io.Reader) error {
	i.Mutex.Lock()
	defer i.Mutex.Unlock()

	if i.src != nil {
		return fmt.Errorf("Maximum number of sources reached: 1")
	}
	i.src = src
	return nil
}
