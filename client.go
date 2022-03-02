package helheim_go

import (
	"fmt"
	"sync"
)

type Client interface {
	NewSession(options CreateSessionOptions) (Session, error)
	GetBalance() (*BalanceResponse, error)
	Version() (*VersionResponse, error)
	GetHelheim() Helheim
}

type client struct {
	logger  Logger
	helheim Helheim
}

var clientContainer = struct {
	sync.Mutex
	instance Client
}{}

func ProvideClient(apiKey string, discover bool, withAutoReAuth bool, logger Logger) (Client, error) {
	clientContainer.Lock()
	defer clientContainer.Unlock()

	if clientContainer.instance != nil {
		return clientContainer.instance, nil
	}

	instance, err := NewClient(apiKey, discover, withAutoReAuth, logger)

	if err != nil {
		return nil, err
	}

	clientContainer.instance = instance

	return clientContainer.instance, nil
}

func NewClient(apiKey string, discover bool, withAutoReAuth bool, logger Logger) (Client, error) {
	if logger == nil {
		logger = NewNoopLogger()
	}

	h, err := newHelheim(apiKey, discover, withAutoReAuth, logger)

	if err != nil {
		return nil, err
	}

	return &client{
		logger:  logger,
		helheim: h,
	}, nil
}

func (c *client) NewSession(options CreateSessionOptions) (Session, error) {
	return newSession(c.logger, c.helheim, options)
}

func (c *client) GetBalance() (*BalanceResponse, error) {
	return c.helheim.GetBalance()
}

func (c *client) Version() (*VersionResponse, error) {
	return c.helheim.Version()
}

func (c *client) GetHelheim() Helheim {
	return c.helheim
}

func (c *client) DeleteSession(sessionId int) error {
	resp, err := c.helheim.DeleteSession(sessionId)

	if err != nil {
		return err
	}

	if resp != nil && resp.Error != false {
		return fmt.Errorf("failed to delete session %d", sessionId)
	}

	return nil
}
