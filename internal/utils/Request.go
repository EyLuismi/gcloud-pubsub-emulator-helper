package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

type ClientInterface interface {
	Get(path string) (Response, error)
	Post(path string, body []byte) (Response, error)
	Put(path string, body []byte) (Response, error)
	Delete(path string) (Response, error)
	Patch(path string, body []byte) (Response, error)
}

type Client struct {
	host    string
	client  http.Client
	version string
}

type Headers = map[string]string

type Response struct {
	StatusCode int
	Body       []byte
}

func (c Client) getUrl(path string) string {
	return fmt.Sprintf("http://%s/%s/%s", c.host, c.version, path)
}

func NewClient(baseUrl string, version string) Client {
	return Client{
		host:    baseUrl,
		client:  http.Client{},
		version: "v1",
	}
}

func (c Client) makeCall(
	method string,
	url string,
	body []byte,
) (Response, error) {

	fmt.Printf("%s %s\n", method, url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}, nil
}

func (c Client) Get(
	path string,
) (Response, error) {
	return c.makeCall(
		http.MethodGet,
		c.getUrl(path),
		[]byte(""),
	)
}

func (c Client) Post(
	path string,
	body []byte,
) (Response, error) {
	return c.makeCall(
		http.MethodPost,
		c.getUrl(path),
		body,
	)
}

func (c Client) Put(
	path string,
	body []byte,
) (Response, error) {
	return c.makeCall(
		http.MethodPut,
		c.getUrl(path),
		body,
	)
}

func (c Client) Delete(
	path string,
) (Response, error) {
	return c.makeCall(
		http.MethodDelete,
		c.getUrl(path),
		[]byte(""),
	)
}

func (c Client) Patch(
	path string,
	body []byte,
) (Response, error) {
	return c.makeCall(
		http.MethodPatch,
		c.getUrl(path),
		body,
	)
}

func IsValidHost(host string) bool {
	hostOnly, port, err := net.SplitHostPort(host)
	if err != nil {
		hostOnly = host
		port = ""
	}

	if hostOnly == "" {
		hostOnly = "localhost"
	}

	if net.ParseIP(hostOnly) != nil {
		return true
	}

	if _, err := net.LookupHost(hostOnly); err != nil {
		return false // Invalid hostname
	}

	if port != "" {
		if p, err := strconv.Atoi(port); err != nil || p < 1 || p > 65535 {
			return false // Invalid Port
		}
	}

	return true
}
