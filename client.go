package helheim_go

import (
	"fmt"
	"sync"
)

type Client interface {
	NewSession(options CreateSessionOptions) (Session, error)
	GetBalance() (*BalanceResponse, error)
	GetHelheim() Helheim
	SetLogger(logger Logger)
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
		logger.Error("failed to create helheim client: %w", err)
		return nil, err
	}

	logger.Info("created new helheim client")

	return &client{
		logger:  logger,
		helheim: h,
	}, nil
}

func (c *client) NewSession(options CreateSessionOptions) (Session, error) {
	s, err := newSession(c.logger, c.helheim, options)

	if err != nil {
		c.logger.Error("failed to create session: %w", err)
		return nil, err
	}

	c.logger.Info("created new session with id: %d", s.GetSessionId())

	return s, nil
}

func (c *client) GetBalance() (*BalanceResponse, error) {
	b, err := c.helheim.GetBalance()

	if err != nil {
		c.logger.Error("failed to retrieve balance: %w", err)
		return nil, err
	}

	return b, nil
}

func (c *client) SetLogger(logger Logger) {
	c.logger = logger
	c.helheim.SetLogger(logger)
}

func (c *client) GetHelheim() Helheim {
	return c.helheim
}

func (c *client) DeleteSession(sessionId int) error {
	resp, err := c.helheim.DeleteSession(sessionId)

	if err != nil {
		c.logger.Error("failed to delete session: %w", err)
		return err
	}

	if resp != nil && resp.Error != false {
		return fmt.Errorf("failed to delete session %d", sessionId)
	}

	return nil
}
