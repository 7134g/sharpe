package execute

import (
	"io/ioutil"
	"net/http"
	"time"
)

func get(u string, header map[string]string) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
