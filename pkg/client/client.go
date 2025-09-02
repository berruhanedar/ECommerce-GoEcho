package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type Client struct {
	http    *fasthttp.Client
	baseURL string
	timeout time.Duration
}

func New(baseURL string, timeout time.Duration) *Client {
	return &Client{
		http:    &fasthttp.Client{},
		baseURL: baseURL,
		timeout: timeout,
	}
}

func (c *Client) Get(path string, headers map[string]string, out interface{}) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.baseURL + path)
	req.Header.SetMethod(fasthttp.MethodGet)
	setHeaders(req, headers)

	if err := c.http.DoTimeout(req, resp, c.timeout); err != nil {
		return err
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK {
		return fmt.Errorf("GET failed: %d", sc)
	}
	return json.Unmarshal(resp.Body(), out)
}

func (c *Client) Post(path string, headers map[string]string, body interface{}, out interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.baseURL + path)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	setHeaders(req, headers)
	req.SetBody(b)

	if err := c.http.DoTimeout(req, resp, c.timeout); err != nil {
		return err
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK && sc != fasthttp.StatusCreated {
		return fmt.Errorf("POST failed: %d", sc)
	}
	return json.Unmarshal(resp.Body(), out)
}

func (c *Client) Put(path string, headers map[string]string, body interface{}, out interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.baseURL + path)
	req.Header.SetMethod(fasthttp.MethodPut)
	req.Header.Set("Content-Type", "application/json")
	setHeaders(req, headers)
	req.SetBody(b)

	if err := c.http.DoTimeout(req, resp, c.timeout); err != nil {
		return err
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK && sc != fasthttp.StatusNoContent {
		return fmt.Errorf("PUT failed: %d", sc)
	}
	if out != nil && len(resp.Body()) > 0 {
		return json.Unmarshal(resp.Body(), out)
	}
	return nil
}

func (c *Client) Patch(path string, headers map[string]string, body interface{}, out interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.baseURL + path)
	req.Header.SetMethod(fasthttp.MethodPatch)
	req.Header.Set("Content-Type", "application/json")
	setHeaders(req, headers)
	req.SetBody(b)

	if err := c.http.DoTimeout(req, resp, c.timeout); err != nil {
		return err
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK && sc != fasthttp.StatusNoContent {
		return fmt.Errorf("PATCH failed: %d", sc)
	}
	if out != nil && len(resp.Body()) > 0 {
		return json.Unmarshal(resp.Body(), out)
	}
	return nil
}

func (c *Client) Delete(path string, headers map[string]string, out interface{}) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c.baseURL + path)
	req.Header.SetMethod(fasthttp.MethodDelete)
	setHeaders(req, headers)

	if err := c.http.DoTimeout(req, resp, c.timeout); err != nil {
		return err
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK && sc != fasthttp.StatusNoContent {
		return fmt.Errorf("DELETE failed: %d", sc)
	}
	if out != nil && len(resp.Body()) > 0 {
		return json.Unmarshal(resp.Body(), out)
	}
	return nil
}

func setHeaders(req *fasthttp.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
