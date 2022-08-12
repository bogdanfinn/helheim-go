package helheim_go

/*
char *auth(char apiKey[], int discover);
char *getBalance();
char *bifrost(int sessionID, char libraryPath[]);
char *wokou(int sessionID, char browser[]);

char *createSession(char options[]);
char *deleteSession(int sessionID);
char *debug(int sessionID, int state);

char *request(int sessionID, char payload[]);

char *setProxy(int sessionID, char proxy[]);
char *setHeaders(int sessionID, char headers[]);
char *setKasada(int sessionID, char kasada[]);
char *setKasadaHooks(int sessionID, char kasadaHooks[]);

char *setCookie(int sessionID, char cookie[]);
char *delCookie(int sessionID, char cookie[]);
#include "Python.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const authValidMinutes = 30

type Helheim interface {
	Auth() (*AuthResponse, error)
	GetBalance() (*BalanceResponse, error)
	CreateSession(options CreateSessionOptions) (*SessionResponse, error)
	DeleteSession(sessionId int) (*SessionDeleteResponse, error)
	Debug(sessionId int, state int) (interface{}, error)
	Request(sessionId int, options RequestOptions) (*RequestResponse, error)
	Wokou(sessionId int, browser string) (*WokouResponse, error)
	SetProxy(sessionId int, proxy string) (*SetProxyResponse, error)
	SetHeaders(sessionId int, headers map[string]string) (*SetHeadersResponse, error)
	SetCookie(sessionId int, cookie SessionCookie) (*ModifyCookiesResponse, error)
	DelCookie(sessionId int, cookieName string) (*ModifyCookiesResponse, error)
	SetKasada(sessionId int, options KasadaOptions) (interface{}, error)
	SetKasadaHooks(sessionId int, options KasadaHooksOptions) (interface{}, error)
	SetLogger(logger Logger)
}

type helheim struct {
	logger         Logger
	apiKey         string
	authLck        sync.Mutex
	lastAuth       *time.Time
	discover       bool
	withAutoReAuth bool
}

func newHelheim(apiKey string, discover bool, withAutoReAuth bool, logger Logger) (Helheim, error) {
	if logger == nil {
		logger = NewNoopLogger()
	}

	h := &helheim{
		logger:         logger,
		apiKey:         apiKey,
		withAutoReAuth: withAutoReAuth,
		discover:       discover,
		authLck:        sync.Mutex{},
	}

	auth, err := h.Auth()

	if err != nil {
		return nil, err
	}

	if auth.Response != "authenticated" {
		return nil, fmt.Errorf("could not authenticate against helheim")
	}

	logger.Info("initiated helheim")

	return h, nil
}

func (h *helheim) Auth() (*AuthResponse, error) {
	h.authLck.Lock()
	defer h.authLck.Unlock()

	if !h.needReAuth() {
		return nil, nil
	}

	apiKey := C.CString(h.apiKey)
	discover := 0

	if h.discover {
		discover = 1
	}

	d := C.int(discover)
	authResp := C.auth(apiKey, d)
	jsonPayload := C.GoString(authResp)

	C.free(unsafe.Pointer(apiKey))

	h.logger.Debug("helheim response for Auth: %s", jsonPayload)
	authResponse := AuthResponse{}
	err := h.handleResponse(jsonPayload, &authResponse)

	now := time.Now()
	h.lastAuth = &now

	return &authResponse, err
}

func (h *helheim) CreateSession(options CreateSessionOptions) (*SessionResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	optionsString, err := json.Marshal(options)

	if err != nil {
		return nil, err
	}

	opt := C.CString(string(optionsString))
	jsonPayload := C.GoString(C.createSession(opt))

	C.free(unsafe.Pointer(opt))

	h.logger.Debug("helheim response for CreateSession: %s", jsonPayload)
	sessionResponse := SessionResponse{}
	err = h.handleResponse(jsonPayload, &sessionResponse)

	return &sessionResponse, err
}

func (h *helheim) GetBalance() (*BalanceResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	jsonPayload := C.GoString(C.getBalance())

	h.logger.Debug("helheim response for GetBalance: %s", jsonPayload)
	balanceResponse := BalanceResponse{}
	err = h.handleResponse(jsonPayload, &balanceResponse)

	return &balanceResponse, err
}

func (h *helheim) DeleteSession(sessionId int) (*SessionDeleteResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.deleteSession(sId))

	h.logger.Debug("helheim response for DeleteSession: %s", jsonPayload)
	deleteResponse := SessionDeleteResponse{}
	err = h.handleResponse(jsonPayload, &deleteResponse)

	return &deleteResponse, err
}

func (h *helheim) Request(sessionId int, options RequestOptions) (*RequestResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	ops := options.Options

	requestOptionsInternal := requestOptionsInternal{
		Method:  options.Method,
		Url:     options.Url,
		Options: ops,
	}

	optionsString, err := json.Marshal(requestOptionsInternal)

	if err != nil {
		return nil, err
	}

	opt := C.CString(string(optionsString))
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.request(sId, opt))
	C.free(unsafe.Pointer(opt))

	requestResponse := RequestResponse{}

	h.logger.Debug("helheim response for Request: %s", jsonPayload)
	err = h.handleResponse(jsonPayload, &requestResponse)

	return &requestResponse, err
}

func (h *helheim) Wokou(sessionId int, browser string) (*WokouResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	b := C.CString(browser)
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.wokou(sId, b))

	C.free(unsafe.Pointer(b))

	h.logger.Debug("helheim response for Wokou: %s", jsonPayload)
	wokouResponse := WokouResponse{}
	err = h.handleResponse(jsonPayload, &wokouResponse)

	return &wokouResponse, err
}

func (h *helheim) SetProxy(sessionId int, proxy string) (*SetProxyResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	p := C.CString(proxy)
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.setProxy(sId, p))

	C.free(unsafe.Pointer(p))

	h.logger.Debug("helheim response for SetProxy: %s", jsonPayload)
	setProxyResponse := SetProxyResponse{}
	err = h.handleResponse(jsonPayload, &setProxyResponse)

	return &setProxyResponse, err
}

func (h *helheim) SetHeaders(sessionId int, headers map[string]string) (*SetHeadersResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	headersString, err := json.Marshal(headers)

	if err != nil {
		return nil, err
	}

	headersParam := C.CString(string(headersString))
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.setHeaders(sId, headersParam))

	C.free(unsafe.Pointer(headersParam))

	h.logger.Debug("helheim response for SetHeaders: %s", jsonPayload)
	setHeadersResponse := SetHeadersResponse{}
	err = h.handleResponse(jsonPayload, &setHeadersResponse)

	return &setHeadersResponse, err
}

func (h *helheim) SetCookie(sessionId int, cookie SessionCookie) (*ModifyCookiesResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	cookiePayload, err := json.Marshal(cookie)

	if err != nil {
		return nil, err
	}

	c := C.CString(string(cookiePayload))
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.setCookie(sId, c))

	h.logger.Debug("helheim response for SetCookie: %s", jsonPayload)
	C.free(unsafe.Pointer(c))

	setCookiesResponse := ModifyCookiesResponse{}
	err = h.handleResponse(jsonPayload, &setCookiesResponse)

	return &setCookiesResponse, err
}

func (h *helheim) DelCookie(sessionId int, cookieName string) (*ModifyCookiesResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	c := C.CString(cookieName)
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.delCookie(sId, c))
	h.logger.Debug("helheim response for DelCookie: %s", jsonPayload)

	C.free(unsafe.Pointer(c))

	setCookiesResponse := ModifyCookiesResponse{}
	err = h.handleResponse(jsonPayload, &setCookiesResponse)

	return &setCookiesResponse, err
}

func (h *helheim) Debug(sessionId int, state int) (interface{}, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	sId := C.int(sessionId)
	stateInt := C.int(state)

	jsonPayload := C.GoString(C.debug(sId, stateInt))
	h.logger.Debug("helheim response for Debug: %s", jsonPayload)

	return jsonPayload, nil
}

func (h *helheim) SetKasada(sessionId int, options KasadaOptions) (interface{}, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	optionsString, err := json.Marshal(options)

	if err != nil {
		return nil, err
	}

	opt := C.CString(string(optionsString))
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.setKasada(sId, opt))
	h.logger.Debug("helheim response for SetKasada: %s", jsonPayload)

	C.free(unsafe.Pointer(opt))

	return jsonPayload, err
}

func (h *helheim) SetKasadaHooks(sessionId int, options KasadaHooksOptions) (interface{}, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	optionsString, err := json.Marshal(options)

	if err != nil {
		return nil, err
	}

	opt := C.CString(string(optionsString))
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.setKasadaHooks(sId, opt))

	C.free(unsafe.Pointer(opt))
	h.logger.Debug("helheim response for SetKasadaHooks: %s", jsonPayload)

	return jsonPayload, err
}

func (h *helheim) SetLogger(logger Logger) {
	h.logger = logger
}

func (h *helheim) handleResponse(jsonPayload string, ret interface{}) error {
	errorResponse := ErrorAwareResponse{}
	err := json.Unmarshal([]byte(jsonPayload), &errorResponse)

	if err != nil {
		e := fmt.Errorf("could not unmarshall helheim response to error aware response: %w", err)
		h.logger.Error("error while unmarshalling helheim response: %w", e)

		return e
	}

	if errorResponse.Error {
		e := fmt.Errorf("helheim error: %s", errorResponse.ErrorMsg)
		h.logger.Error("error received in helheim response: %w", e)

		return e
	}

	err = json.Unmarshal([]byte(jsonPayload), ret)

	if err != nil {
		e := fmt.Errorf("could not unmarshall helheim response to response type: %w", err)
		h.logger.Error("error while unmarshalling helheim response: %w", e)
		return e
	}

	return nil
}

func (h *helheim) needReAuth() bool {
	if h.lastAuth == nil {
		return true
	}

	now := time.Now()

	minutes := now.Sub(*h.lastAuth).Minutes()

	h.logger.Info("%d minutes since last authentication", minutes)

	needsReAuth := minutes >= authValidMinutes

	h.logger.Debug("%d minutes since last auth. need re auth: %v", minutes, needsReAuth)

	return needsReAuth
}

func (h *helheim) reAuth() error {
	if !h.withAutoReAuth {
		h.logger.Debug("auto re auth not enabled. skipping")
		return nil
	}

	if !h.needReAuth() {
		return nil
	}

	_, err := h.Auth()

	if err != nil {
		h.logger.Error("failed to authenticate helheim: %w", err)
		return err
	}

	return nil
}
