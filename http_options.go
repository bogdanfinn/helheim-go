package helheim_go

type HttpClientOption func(config *httpClientConfig)

type httpClientConfig struct {
	proxyUrl         string
	debug            bool
	wokouBrowser     string
	isBase64Response bool
}

func WithProxyUrl(proxyUrl string) HttpClientOption {
	return func(config *httpClientConfig) {
		config.proxyUrl = proxyUrl
	}
}

func WithBase64Response() HttpClientOption {
	return func(config *httpClientConfig) {
		config.isBase64Response = true
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
