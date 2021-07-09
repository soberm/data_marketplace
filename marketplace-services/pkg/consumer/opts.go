package consumer

import (
	"encoding/json"
	"os"
)

type options struct {
	ConfigFile    string
	ProxyConfig   ProxyConfig   `json:"proxyConfig"`
	LoggingConfig LoggingConfig `json:"loggingConfig"`
	SearchConfig  SearchConfig  `json:"searchConfig"`
}

type ProxyConfig struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Account  string `json:"account"`
}

type LoggingConfig struct {
	Verbosity int `json:"verbosity"`
}

type SearchConfig struct {
	BrokerLocation int    `json:"brokerLocation"`
	DataType       string `json:"dataType"`
	MinCost        uint64 `json:"minCost"`
	MaxCost        uint64 `json:"maxCost"`
	MinFrequency   uint64 `json:"minFrequency"`
	MaxFrequency   uint64 `json:"maxFrequency"`
}

type Option interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

func (funcOpt *funcOption) apply(options *options) {
	funcOpt.f(options)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func defaultOptions() options {
	return options{
		ConfigFile: "./configs/consumer/config.json",
		ProxyConfig: ProxyConfig{
			Address:  "127.0.0.1",
			Port:     25566,
			Username: "michael",
			Password: "12345678",
			Account:  "0x6da49C19d815c1c61050046456398599720716A2",
		},
		LoggingConfig: LoggingConfig{
			Verbosity: 4,
		},
		SearchConfig: SearchConfig{
			BrokerLocation: 6,
			DataType:       "light",
			MinCost:        0,
			MaxCost:        20,
			MinFrequency:   1,
			MaxFrequency:   10,
		},
	}
}

func (o *options) loadConfiguration() error {
	configFile, err := os.Open(o.ConfigFile)
	if err != nil {
		return err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	return jsonParser.Decode(&o)
}

func WithConfigFile(configFile string) Option {
	return newFuncOption(func(o *options) {
		o.ConfigFile = configFile
	})
}

func WithProxyConfig(proxyConfig ProxyConfig) Option {
	return newFuncOption(func(o *options) {
		o.ProxyConfig = proxyConfig
	})
}

func WithLoggingConfig(loggingConfig LoggingConfig) Option {
	return newFuncOption(func(o *options) {
		o.LoggingConfig = loggingConfig
	})
}

func WithSearchConfig(searchConfig SearchConfig) Option {
	return newFuncOption(func(o *options) {
		o.SearchConfig = searchConfig
	})
}
