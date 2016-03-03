package buildkite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	orgURL   *url.URL
	apiToken string
}

func NewClient(orgURLStr, apiToken string) (*Client, error) {
	orgURL, err := url.Parse(orgURLStr)
	if err != err {
		return nil, err
	}

	return &Client{
		orgURL:   orgURL,
		apiToken: apiToken,
	}, nil
}

func (c *Client) Get(pathParts []string, resBody interface{}) error {
	return c.doJSON("GET", pathParts, nil, resBody)
}

func (c *Client) createRawRequest(method string, pathParts []string, reqBodyBytes []byte) *http.Request {
	urlPath := &url.URL{
		Path: strings.Join(pathParts, "/"),
	}
	reqURL := c.orgURL.ResolveReference(urlPath)

	req := &http.Request{
		Method: method,
		Header: http.Header{},
		URL:    reqURL,
	}
	req.Header.Add("Authorization", "Bearer "+c.apiToken)

	return req
}

func (c *Client) doRaw(req *http.Request) ([]byte, error) {
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%s", res.Status)
	}

	return resBodyBytes, nil
}

func (c *Client) doJSON(method string, pathParts []string, reqBody, resBody interface{}) error {
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req := c.createRawRequest(method, pathParts, reqBodyBytes)

	resBodyBytes, err := c.doRaw(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(resBodyBytes, resBody)
}
