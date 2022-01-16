package helheim_go

import (
	"fmt"
	"net/http"
	"time"
)

type Session interface {
	Delete() error
	Debug(state int) (interface{}, error)
	Request(options RequestOptions) (*RequestResponse, error)
	Bifrost(libraryPath string) (interface{}, error)
	Wokou(browser string) (*WokouResponse, error)
	SetProxy(proxy string) (*SetProxyResponse, error)
	SetHeaders(headers map[string]string) (*SetHeadersResponse, error)
	SetCookie(cookie string) (interface{}, error)
	DelCookie(cookie string) (interface{}, error)
	SetKasada(options KasadaOptions) (interface{}, error)
	SetKasadaHooks(options KasadaHooksOptions) (interface{}, error)
	GetGoHttpCookies() []*http.Cookie
	GetSessionId() int
	GetHeaders() map[string]string
	GetCookies() []SessionCookie
}

type session struct {
	logger         Logger
	helheim        Helheim
	helheimSession SessionResponse
	sessionId      int
	headers        map[string]string
	cookies        []SessionCookie
}

func newSession(logger Logger, helheim Helheim, options CreateSessionOptions) (Session, error) {
	helheimSession, err := helheim.CreateSession(options)

	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = NewNoopLogger()
	}

	return &session{
		logger:         logger,
		helheim:        helheim,
		helheimSession: *helheimSession,
		sessionId:      helheimSession.SessionId,
		headers:        helheimSession.Headers,
		cookies:        helheimSession.Cookies,
	}, nil
}

func (s *session) Request(options RequestOptions) (*RequestResponse, error) {
	resp, err := s.helheim.Request(s.GetSessionId(), options)

	if err != nil {
		return nil, err
	}

	for key, value := range resp.Session.Headers {
		s.headers[key] = value
	}

	s.cookies = append(s.cookies, resp.Session.Cookies...)

	return resp, nil
}

func (s *session) Bifrost(libraryPath string) (interface{}, error) {
	return s.helheim.Bifrost(s.GetSessionId(), libraryPath)
}

func (s *session) Wokou(browser string) (*WokouResponse, error) {
	return s.helheim.Wokou(s.GetSessionId(), browser)
}

func (s *session) SetProxy(proxy string) (*SetProxyResponse, error) {
	return s.helheim.SetProxy(s.GetSessionId(), proxy)
}

func (s *session) SetHeaders(headers map[string]string) (*SetHeadersResponse, error) {
	for key, value := range headers {
		s.headers[key] = value
	}

	return s.helheim.SetHeaders(s.GetSessionId(), headers)
}

func (s *session) SetCookie(cookie string) (interface{}, error) {
	return s.helheim.SetCookie(s.GetSessionId(), cookie)
}

func (s *session) DelCookie(cookie string) (interface{}, error) {
	return s.helheim.DelCookie(s.GetSessionId(), cookie)
}

func (s *session) Debug(state int) (interface{}, error) {
	return s.helheim.Debug(s.GetSessionId(), state)
}

func (s *session) SetKasada(options KasadaOptions) (interface{}, error) {
	return s.helheim.SetKasada(s.GetSessionId(), options)
}

func (s *session) SetKasadaHooks(options KasadaHooksOptions) (interface{}, error) {
	return s.helheim.SetKasadaHooks(s.GetSessionId(), options)
}

func (s *session) GetSessionId() int {
	return s.sessionId
}

func (s *session) GetHeaders() map[string]string {
	return s.headers
}

func (s *session) GetGoHttpCookies() []*http.Cookie {
	var goCookies []*http.Cookie

	for _, cookie := range s.GetCookies() {
		goCookies = append(goCookies, &http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Path:    cookie.Path,
			Domain:  cookie.Domain,
			Expires: time.Unix(int64(cookie.Expires), 0),
		})
	}

	return goCookies
}

func (s *session) GetCookies() []SessionCookie {
	return s.cookies
}

func (s *session) Delete() error {
	resp, err := s.helheim.DeleteSession(s.GetSessionId())

	if err != nil {
		return err
	}

	if resp != nil && resp.Error != false {
		return fmt.Errorf("failed to delete session %d", s.GetSessionId())
	}

	println("")
	//s.logger.Infof("deleted session with id %d", s.GetSessionId())

	return nil
}
