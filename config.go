package line

// Config is the configuration of the connector.
type Config struct {
	// Address is the address of the NATS server.
	Address string `json:"address"`
	// Token is the token to use to connect to the NATS server.
	Token string `json:"token"`
	// Name is the name of the connector used for monitoring and debugging purposes.
	Name string `json:"name"`
}

// DefaultConfig returns a configuration with the default values filled out.
func DefaultConfig() Config {
	c := Config{}
	c.Address = "localhost:4222"
	return c
}
