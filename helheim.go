package helheim_go

/*
char *auth(char apiKey[], int discover);
char *getBalance();
char *helheimVersion();
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
	Version() (*VersionResponse, error)
	CreateSession(options CreateSessionOptions) (*SessionResponse, error)
	DeleteSession(sessionId int) (*SessionDeleteResponse, error)
	Debug(sessionId int, state int) (interface{}, error)
	Request(sessionId int, options RequestOptions) (*RequestResponse, error)
	Bifrost(sessionId int, libraryPath string) (interface{}, error)
	Wokou(sessionId int, browser string) (*WokouResponse, error)
	SetProxy(sessionId int, proxy string) (*SetProxyResponse, error)
	SetHeaders(sessionId int, headers map[string]string) (*SetHeadersResponse, error)
	SetCookie(sessionId int, cookie string) (interface{}, error)
	DelCookie(sessionId int, cookie string) (interface{}, error)
	SetKasada(sessionId int, options KasadaOptions) (interface{}, error)
	SetKasadaHooks(sessionId int, options KasadaHooksOptions) (interface{}, error)
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

	return h, err
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

	balanceResponse := BalanceResponse{}
	err = h.handleResponse(jsonPayload, &balanceResponse)

	return &balanceResponse, err
}

func (h *helheim) Version() (*VersionResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	jsonPayload := C.GoString(C.helheimVersion())

	versionResponse := VersionResponse{}
	err = h.handleResponse(jsonPayload, &versionResponse)

	return &versionResponse, err
}

func (h *helheim) DeleteSession(sessionId int) (*SessionDeleteResponse, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.deleteSession(sId))

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
	bodyValue, ok := ops["body"]

	if !ok || bodyValue == "" {
		if options.Body != "" {
			ops["body"] = options.Body
		}
	}

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

	err = h.handleResponse(jsonPayload, &requestResponse)

	return &requestResponse, err
}

func (h *helheim) Bifrost(sessionId int, libraryPath string) (interface{}, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	lp := C.CString(libraryPath)
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.bifrost(sId, lp))

	C.free(unsafe.Pointer(lp))

	return jsonPayload, nil
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

	setHeadersResponse := SetHeadersResponse{}
	err = h.handleResponse(jsonPayload, &setHeadersResponse)

	return &setHeadersResponse, err
}

func (h *helheim) SetCookie(sessionId int, cookie string) (interface{}, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	c := C.CString(cookie)
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.setCookie(sId, c))

	C.free(unsafe.Pointer(c))

	return jsonPayload, nil
}

func (h *helheim) DelCookie(sessionId int, cookie string) (interface{}, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	c := C.CString(cookie)
	sId := C.int(sessionId)

	jsonPayload := C.GoString(C.delCookie(sId, c))

	C.free(unsafe.Pointer(c))

	return jsonPayload, nil
}

func (h *helheim) Debug(sessionId int, state int) (interface{}, error) {
	err := h.reAuth()
	if err != nil {
		return nil, err
	}

	sId := C.int(sessionId)
	stateInt := C.int(state)

	jsonPayload := C.GoString(C.debug(sId, stateInt))

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

	return jsonPayload, err
}

func (h *helheim) handleResponse(jsonPayload string, ret interface{}) error {
	errorResponse := ErrorAwareResponse{}
	err := json.Unmarshal([]byte(jsonPayload), &errorResponse)

	if err != nil {
		return err
	}

	if errorResponse.Error {
		return fmt.Errorf("error: %s", errorResponse.ErrorMsg)
	}

	err = json.Unmarshal([]byte(jsonPayload), ret)

	return err
}

func (h *helheim) needReAuth() bool {
	if h.lastAuth == nil {
		return true
	}

	now := time.Now()

	minutes := now.Sub(*h.lastAuth).Minutes()

	h.logger.Info("checked reauth minutes", minutes)

	return minutes >= authValidMinutes
}

func (h *helheim) reAuth() error {
	if !h.withAutoReAuth {
		return nil
	}

	if !h.needReAuth() {
		return nil
	}

	_, err := h.Auth()

	return err
}
