package helheim_go

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	GetSessionHeaders() map[string]string
	GetSessionCookies() []SessionCookie
}

type httpClient struct {
	logger  Logger
	config  *httpClientConfig
	session Session
}

func newHttpClient(logger Logger, session Session, options ...HttpClientOption) HttpClient {
	config := &httpClientConfig{}
	for _, opt := range options {
		opt(config)
	}

	return &httpClient{
		logger:  logger,
		session: session,
		config:  config,
	}
}

func (c *httpClient) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}
func (c *httpClient) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}
func (c *httpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return c.Do(req)
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	if c.config.debug {
		_, err := c.session.Debug(1)

		if err != nil {
			c.logger.Error("failed to set wokou on http client default session: %w", err)
			return nil, err
		}
	}

	if c.config.bifrostLibraryPath != "" {
		_, err := c.session.Bifrost(c.config.bifrostLibraryPath)

		if err != nil {
			c.logger.Error("failed to set bifrost library path on http client default session: %w", err)
			return nil, err
		}
	}

	if c.config.wokouBrowser != "" {
		wokouResp, err := c.session.Wokou(c.config.wokouBrowser)

		if err != nil {
			c.logger.Error("failed to set wokou on http client default session: %w", err)
			return nil, err
		}

		if wokouResp.Error {
			err = fmt.Errorf("received error response from helheim: %s", wokouResp.ErrorMsg)
			c.logger.Error("failed to set wokou on http client default session: %w", err)

			return nil, err
		}
	}

	if len(req.Header) > 0 {
		headerMap := make(map[string]string, 0)
		for key, value := range req.Header {
			headerMap[key] = value[0] // TODO: find a better handling here
		}

		headerResp, err := c.session.SetHeaders(headerMap)

		if err != nil {
			c.logger.Error("failed to set header on http client default session: %w", err)
			return nil, err
		}

		if headerResp.Error {
			err = fmt.Errorf("received error response from helheim: %s", headerResp.ErrorMsg)
			c.logger.Error("failed to set header on http client default session: %w", err)

			return nil, err
		}
	}

	if c.config.proxyUrl != "" {
		proxyResp, err := c.session.SetProxy(c.config.proxyUrl)

		if err != nil {
			c.logger.Error("failed to set proxy on http client default session: %w", err)
			return nil, err
		}

		if proxyResp.Error {
			err = fmt.Errorf("received error response from helheim: %s", proxyResp.ErrorMsg)
			c.logger.Error("failed to set proxy on http client default session: %w", err)

			return nil, err
		}
	}

	opts := make(map[string]string, 0)

	var body string
	if req.Body != nil {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			c.logger.Error("failed to prepare helheim request body: %w", err)

			return nil, err
		}
		body = string(bodyBytes)

		if req.Header.Get("Content-Type") == "application/json" {
			opts["json"] = body
		} else {
			opts["data"] = body
		}
	}

	reqOpts := RequestOptions{
		Method:  req.Method,
		Url:     req.URL.String(),
		Options: opts,
	}

	resp, err := c.session.Request(reqOpts)

	if err != nil {
		c.logger.Error("failed to get response on http client default session: %w", err)
		return nil, err
	}

	if resp.Error {
		err = fmt.Errorf("received error response from helheim: %s", resp.ErrorMsg)
		c.logger.Error("failed to get response on http client default session: %w", err)

		return nil, err
	}

	return &http.Response{
		Status:     fmt.Sprintf("Status Code: %d", resp.Response.StatusCode),
		StatusCode: resp.Response.StatusCode,
		//Proto:            "",
		//ProtoMajor:       0,
		//ProtoMinor:       0,
		Header:        toGoHeader(resp.Response.Headers),
		Body:          io.NopCloser(strings.NewReader(resp.Response.Body)),
		ContentLength: int64(len(resp.Response.Body)), // TODO: is this correct
		//TransferEncoding: nil,
		//Close:            false,
		//Uncompressed:     false,
		//Trailer:          nil,
		Request: req,
		//TLS:              nil,
	}, nil
}

func (c *httpClient) GetSessionHeaders() map[string]string {
	return c.session.GetHeaders()
}

func (c *httpClient) GetSessionCookies() []SessionCookie {
	return c.session.GetCookies()
}

func toGoHeader(hdrMap map[string]string) http.Header {
	hdr := http.Header{}

	for key, value := range hdrMap {
		hdr.Add(key, value)
	}

	return hdr
}
