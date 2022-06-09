package helheim_go

type HttpClientOption func(config *httpClientConfig)

type httpClientConfig struct {
	proxyUrl     string
	debug        bool
	wokouBrowser string
}

func WithProxyUrl(proxyUrl string) HttpClientOption {
	return func(config *httpClientConfig) {
		config.proxyUrl = proxyUrl
	}
}

func WithWokou(browser string) HttpClientOption {
	return func(config *httpClientConfig) {
		config.wokouBrowser = browser
	}
}

func WithDebug() HttpClientOption {
	return func(config *httpClientConfig) {
		config.debug = true
	}
}
