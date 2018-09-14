package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	_ "strings"
)

func HTTPClient() *http.Client {

	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}

	return &http.Client{
		Transport: tr,
	}
}

func Call(method, path, master string, data interface{}) (io.ReadCloser, int, error) {

	params := bytes.NewBuffer(nil)

	if data != nil {
		buf, err := json.Marshal(data)

		if err != nil {
			return nil, -1, err
		}

		if _, err := params.Write(buf); err != nil {
			return nil, -1, err
		}
	}

	req, err := http.NewRequest(method, path, params)

	if err != nil {
		return nil, -1, err
	}
	req.URL.Host = master
	req.URL.Scheme = "http"

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	} else if method == "POST" {
		req.Header.Set("Content-Type", "application/text")
	}

	resp, err := HTTPClient().Do(req)

	if err != nil {
		return nil, -1, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, -1, err
		}

		if len(body) == 0 {
			return nil, resp.StatusCode, fmt.Errorf("Error: request returned %s, check if the server supports the requested API version",
				http.StatusText(resp.StatusCode))
		}

		return nil, resp.StatusCode, fmt.Errorf("Error response form daemon: %s", bytes.TrimSpace(body))
	}

	return resp.Body, resp.StatusCode, nil

}

// func Stream(method, path, master string, stdout io.Writer, header map[string][]string) error {
// 	req, err := http.NewRequest(method, path, nil)
// 	req.URL.Host = master
// 	req.URL.Scheme = "http"
// 	if header != nil {
// 		for k, v := range header {
// 			req.Header[k] = v
// 		}
// 	}

// 	resp, err := HTTPClient(master).Do(req)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return err
// 	}

// 	StdCopy(stdout, resp.Body)
// 	return nil
// }

func Stream(method, path, master string) (io.ReadCloser, int, error) {
	req, err := http.NewRequest(method, path, nil)
	req.URL.Host = master
	req.URL.Scheme = "http"

	resp, err := HTTPClient().Do(req)

	if err != nil {
		return nil, -1, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, -1, err
		}

		if len(body) == 0 {
			return nil, resp.StatusCode, fmt.Errorf("Error: request returned %s, check if the server supports the requested API version",
				http.StatusText(resp.StatusCode))
		}

		return nil, resp.StatusCode, fmt.Errorf("Error response form daemon: %s", bytes.TrimSpace(body))
	}

	return resp.Body, resp.StatusCode, nil
}
