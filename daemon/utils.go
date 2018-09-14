package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"strings"
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

func StdCopy(dstout io.Writer, src io.Reader) error {
	rc := newBufReader(src)
	_, err := io.Copy(dstout, rc)

	return err
}

func readBody(stream io.ReadCloser, statusCode int, err error) ([]byte, int, error) {
	if stream != nil {
		defer stream.Close()
	}

	if err != nil {
		return nil, statusCode, err
	}

	body, err := ioutil.ReadAll(stream)

	if err != nil {
		return nil, -1, err
	}

	return body, statusCode, nil

}
func normalization(data []int, number int) (nordata []float64) { //归一化
	maxdata := float64(findmax(data))
	mindata := float64(findmin(data))
	nordata = make([]float64, 0, number)
	for _, v := range data {
		if maxdata == mindata {
			nordata = append(nordata, 1)
		} else {
			nordata = append(nordata, (float64(v)-mindata)/(maxdata-mindata))
		}
		//nordata[i] = (v - mindata) / (maxdata - mindata)
	}
	return nordata
}
func normalizationfloat(data []float64) (nordata []float64) {
	maxdata := findmaxfloat(data)
	mindata := findminfloat(data)
	for _, v := range data {
		if maxdata == mindata {
			nordata = append(nordata, 1)
		} else {
			nordata = append(nordata, (v-mindata)/(maxdata-mindata))
		}
		//nordata[i] = (v - mindata) / (maxdata - mindata)
	}
	return nordata
}
func findmax(data []int) (maxdata int) {
	maxdata = data[0]
	for _, v := range data {
		if v > maxdata {
			maxdata = v
		}
	}
	return maxdata
}
func findmaxfloat(data []float64) (maxdata float64) {
	maxdata = data[0]
	for _, v := range data {
		if v > maxdata {
			maxdata = v
		}
	}
	return maxdata
}
func findmin(data []int) (mindata int) {
	mindata = data[0]
	for _, v := range data {
		if v < mindata {
			mindata = v
		}
	}
	return mindata
}
func findminfloat(data []float64) (mindata float64) {
	mindata = data[0]
	for _, v := range data {
		if v < mindata {
			mindata = v
		}
	}
	return mindata
}
func comparestring(a, b string) int {
	if strings.Compare(a, b) == 0 {
		return 1
	} else {
		return 0
	}
}
func comparemap(x, y map[string]string) int { //map比较
	if (len(x)) != len(y) {
		return 0
	}
	for k, xv := range x {
		if yv, ok := y[k]; !ok || yv != xv {
			return 0
		}
	}
	return 1
}
func findmaxnode(data []float64) (maxdata float64, index int) {
	maxdata = data[0]
	index = 1
	for i, v := range data {
		if v > maxdata {
			maxdata = v
			index = i + 1
		}
	}
	return maxdata, index
}
func mapequal(x, y map[string]string) bool {
	if (len(x)) != len(y) {
		return false
	}
	for k, xv := range x {
		if yv, ok := y[k]; !ok || yv != xv {
			return false
		}
	}
	return true
}

func (nodes *Nodes) ToJsonString() string { //json转换
	strs, err := json.Marshal(nodes)

	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return string(strs)
}
func (containers *Containers) ToJsonString() string {
	strs, err := json.Marshal(containers)

	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return string(strs)
}
func printSlice(x []float64) {
	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}
func (nodes *Nodes) Len() int {
	return len(nodes.AllNodes)
}
func (nodes *Nodes) Less(i, j int) bool {
	return nodes.AllNodes[i].Index < nodes.AllNodes[j].Index
}
func (nodes *Nodes) Swap(i, j int) {
	nodes.AllNodes[i], nodes.AllNodes[j] = nodes.AllNodes[j], nodes.AllNodes[i]
}
