package helheim_go

type HttpClientOption func(config *httpClientConfig)

type httpClientConfig struct {
	proxyUrl           string
	debug              bool
	wokouBrowser       string
	bifrostLibraryPath string
}

func WithProxyUrl(proxyUrl string) HttpClientOption {
	return func(config *httpClientConfig) {
		config.proxyUrl = proxyUrl
	}
}

func WithBifrost(libraryPath string) HttpClientOption {
	return func(config *httpClientConfig) {
		config.bifrostLibraryPath = libraryPath
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
