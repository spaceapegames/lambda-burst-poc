package client

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	address string
	http.Client
}

func NewClient(address string) Client {
	return Client{
		address: address,
		Client:  http.Client{},
	}
}

func (c Client) Go() (int, string, time.Duration, error) {
	req, err := http.NewRequest(http.MethodGet, c.address+"/doThing", nil)
	if err != nil {
		return 0, "", 0, err
	}
	req.Close = true
	startTime := time.Now()
	resp, err := c.Do(req)
	duration := time.Now().Sub(startTime)
	if err != nil {
		return 0, "", 0, err
	}

	o, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(o), duration, nil
}
