package helheim_go

type Client interface {
	NewSession(options CreateSessionOptions) (Session, error)
	GetBalance() (*BalanceResponse, error)
}

type client struct {
	logger  Logger
	helheim Helheim
}

func NewClient(apiKey string, discover bool, logger Logger) (Client, error) {
	if logger == nil {
		logger = NewNoopLogger()
	}

	h, err := newHelheim(apiKey, discover, logger)

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
