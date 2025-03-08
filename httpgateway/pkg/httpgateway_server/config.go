package httpgateway_server

//nolint:gochecknoglobals
var (
	// DefaultAddress is the default Address to listen to.
	DefaultAddress = ":8080"

	// DefaultConfigSection is the section key used in config files used to
	// configure the gateway options.
	DefaultConfigSection = "httpgateway"
)

// Option is a functional option type for the httpgateway.
type Option func(*Config)

// Config represents the config for the httpgateway.
type Config struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Address string `json:"address" yaml:"address"`
}

// NewConfig creates a config to use with a httpgateway.
func NewConfig() Config {
	return Config{
		Enabled: true,
		Address: DefaultAddress,
	}
}
