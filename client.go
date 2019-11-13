package sfox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	host = "https://api.sfox.com"
)

type Client struct {
	client *http.Client
	apiKey string
	host   string
}

func New(apiKey string) *Client {
	return NewWithHost(apiKey, host)
}

func NewWithHost(apiKey, host string) *Client {
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
		},
		apiKey: apiKey,
		host:   host,
	}
}

func (c *Client) doGet(path string, params url.Values, result interface{}) error {
	url := c.host + path
	query := params.Encode()
	if query != "" {
		url += "?" + query
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return c.do(req, result)
}

func (c *Client) doPost(path string, body interface{}, result interface{}) error {
	url := c.host + path

	buf := bytes.NewBuffer([]byte("{}"))
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	return c.do(req, result)
}

func (c *Client) doDelete(path string) error {
	url := c.host + path

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

func (c *Client) do(req *http.Request, result interface{}) error {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest && res.StatusCode < http.StatusInternalServerError {
		body, _ := ioutil.ReadAll(res.Body)
		return ErrHttpClient{ErrHttp{StatusCode: res.StatusCode, Text: string(body)}}
	} else if res.StatusCode >= http.StatusInternalServerError {
		body, _ := ioutil.ReadAll(res.Body)
		return ErrHttpServer{ErrHttp{StatusCode: res.StatusCode, Text: string(body)}}
	}

	if result != nil {
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return err
		}
	}
	return nil
}
