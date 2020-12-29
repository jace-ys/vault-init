package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type InitRequest struct {
	SecretShares    int `json:"secret_shares"`
	SecretThreshold int `json:"secret_threshold"`
}

type InitResponse struct {
	Keys       []string `json:"keys"`
	KeysBase64 []string `json:"keys_base64"`
	RootToken  string   `json:"root_token"`
}

type UnsealRequest struct {
	Key string `json:"key"`
}

type UnsealResponse struct {
	Sealed   bool `json:"sealed"`
	T        int  `json:"t"`
	N        int  `json:"n"`
	Progress int  `json:"progress"`
}

type Client struct {
	addr   *url.URL
	client *http.Client
}

func NewClient(addr *url.URL) *Client {
	return &Client{
		addr:   addr,
		client: &http.Client{},
	}
}

func (c *Client) Health(ctx context.Context) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodHead, "/v1/sys/health", nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Init(ctx context.Context, opts *InitRequest) (*InitResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPut, "/v1/sys/init", opts)
	if err != nil {
		return nil, err
	}

	resp, body, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d, expected 200", resp.StatusCode)
	}

	data := new(InitResponse)
	err = json.Unmarshal(body, data)

	return data, err
}

func (c *Client) Unseal(ctx context.Context, opts *UnsealRequest) (*UnsealResponse, error) {
	request, err := c.newRequest(ctx, http.MethodPut, "/v1/sys/unseal", opts)
	if err != nil {
		return nil, err
	}

	response, body, err := c.do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d, expected 200", response.StatusCode)
	}

	data := new(UnsealResponse)
	err = json.Unmarshal(body, data)

	return data, err
}

func (c *Client) newRequest(ctx context.Context, method, requestPath string, body interface{}) (*http.Request, error) {
	requestURL := *c.addr
	requestURL.Path = path.Join(requestURL.Path, requestPath)

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	return http.NewRequestWithContext(ctx, method, requestURL.String(), buf)
}

func (c *Client) do(req *http.Request) (*http.Response, []byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, body, nil
}
