package cluster

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type NopFlusher struct{}

func (f *NopFlusher) Flush() {}

type WriterFlusher struct {
	sync.Mutex
	w       io.Writer
	flusher http.Flusher
}

func (wf *WriterFlusher) Write(b []byte) (n int, err error) {
	wf.Lock()
	defer wf.Unlock()
	n, err = wf.w.Write(b)
	wf.flusher.Flush()
	return
}

func (wf *WriterFlusher) Flush() {
	wf.Lock()
	defer wf.Unlock()
	wf.flusher.Flush()
}

func NewWriterFlusher(w io.Writer) *WriterFlusher {
	var flusher http.Flusher
	if f, ok := w.(http.Flusher); ok {
		flusher = f
	} else {
		flusher = &NopFlusher{}
	}

	return &WriterFlusher{w: w, flusher: flusher}
}

type bufReader struct {
	sync.Mutex
	buf    *bytes.Buffer
	reader io.Reader
	err    error
	wait   sync.Cond
}

func newBufReader(r io.Reader) *bufReader {
	reader := &bufReader{
		buf:    &bytes.Buffer{},
		reader: r,
	}

	reader.wait.L = &reader.Mutex
	go reader.drain()
	return reader

}

func (r *bufReader) drain() {
	buf := make([]byte, 1024)
	for {
		n, err := r.reader.Read(buf)
		r.Lock()
		if err != nil {
			r.err = err
		} else {
			r.buf.Write(buf[0:n])
		}
		r.wait.Signal()
		r.Unlock()

		if err != nil {
			break
		}
	}
}

func (r *bufReader) Read(p []byte) (n int, err error) {
	r.Lock()
	defer r.Unlock()

	for {
		n, err := r.buf.Read(p)
		if n > 0 {
			return n, err
		}

		if r.err != nil {
			return 0, r.err
		}

		r.wait.Wait()
	}
}

func (r *bufReader) Close() error {
	closer, ok := r.reader.(io.ReadCloser)

	if !ok {
		return nil
	}

	return closer.Close()
}

func StdCopy(dstout io.Writer, src io.Reader) {
	rc := newBufReader(src)
	io.Copy(dstout, rc)
}

func ReadBody(stream io.ReadCloser, statusCode int, err error) ([]byte, int, error) { //test

	if stream != nil {
		defer stream.Close()
	}
	if err != nil {
		return nil, statusCode, err
	}

	body, err := ioutil.ReadAll(stream)

	if err != nil {
		return nil, statusCode, err
	}

	return body, statusCode, nil

}

func Mkdir(filename string) error {
	if _, _, err := Call("GET", "/create?filename="+filename, "slave1:8090", nil); err != nil {
		return err
	}

	return nil
}

func Remove(filename string) error {
	if _, _, err := Call("GET", "/delete?filename="+filename, "slave1:8090", nil); err != nil {
		return err
	}
	return nil
}
