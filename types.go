package helheim_go

type SessionAwareResponse struct {
	SessionId int `json:"sessionID"`
}

type ErrorAwareResponse struct {
	Error    bool   `json:"error"`
	ErrorMsg string `json:"errorMsg"`
}

type AuthResponse struct {
	SessionAwareResponse
	Response string `json:"response"`
}

type VersionResponse struct {
	ErrorAwareResponse
	SessionAwareResponse
	Version string `json:"version"`
}

type BalanceResponse struct {
	ErrorAwareResponse
	SessionAwareResponse
	Response struct {
		Balance   int  `json:"balance"`
		IsExpired bool `json:"isExpired"`
		Expiry    int  `json:"expiry"`
	} `json:"response"`
}

type SessionResponse struct {
	ErrorAwareResponse
	SessionAwareResponse
	Headers map[string]string `json:"headers"`
	Cookies []SessionCookie   `json:"cookies"`
}

type SessionCookie struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Domain  string `json:"domain"`
	Path    string `json:"path"`
	Expires int    `json:"expires"`
}

type RequestResponse struct {
	ErrorAwareResponse
	SessionAwareResponse
	Session  RequestResponseSession  `json:"session"`
	Response RequestResponseResponse `json:"response"`
}

type RequestResponseSession struct {
	Headers map[string]string `json:"headers"`
	Cookies []SessionCookie   `json:"cookies"`
}

type RequestResponseResponse struct {
	Headers    map[string]string `json:"headers"`
	Cookies    map[string]string `json:"cookies"`
	StatusCode int               `json:"status_code"`
	Body       string            `json:"body"`
}

type SessionDeleteResponse struct {
	ErrorAwareResponse
	SessionAwareResponse
}

type SetHeadersResponse struct {
	ErrorAwareResponse
	SessionAwareResponse
}

type SetProxyResponse struct {
	ErrorAwareResponse
	SessionAwareResponse
}

type WokouResponse struct {
	ErrorAwareResponse
	SessionAwareResponse
	Response string `json:"response"`
}

type CreateSessionOptions struct {
	Browser BrowserOptions `json:"browser"`
	Captcha CaptchaOptions `json:"captcha"`
}

type BrowserOptions struct {
	Browser  string `json:"browser"`
	Mobile   bool   `json:"mobile"`
	Platform string `json:"platform"`
}

type CaptchaOptions struct {
	Provider string `json:"provider"`
}

type RequestOptions struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Body    string            `json:"body"`
	Options map[string]string `json:"options"`
}

type requestOptionsInternal struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Options map[string]string `json:"options"`
}

type KasadaOptions struct {
}

type KasadaHooksOptions struct {
}
